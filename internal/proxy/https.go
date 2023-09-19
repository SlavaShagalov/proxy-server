package proxy

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"github.com/SlavaShagalov/proxy-server/internal/requests"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
)

func (p *Proxy) genCert(host string) (tls.Certificate, error) {
	if _, err := os.Stat(host + ".crt"); os.IsNotExist(err) {
		err = exec.Command("/bin/sh", "gen_cert.sh", host).Run()
		if err != nil {
			p.log.Error("gen_cert.sh failed", zap.Error(err))
			return tls.Certificate{}, err
		}
	}

	cf, err := os.ReadFile(host + ".crt")
	if err != nil {
		p.log.Error("Read .crt file failed", zap.Error(err))
		return tls.Certificate{}, err
	}
	kf, err := os.ReadFile(host + ".key")
	if err != nil {
		p.log.Error("Read .key file failed", zap.Error(err))
		return tls.Certificate{}, err
	}

	tlsCert, err := tls.X509KeyPair(cf, kf)
	if err != nil {
		p.log.Error("Parse public/private key failed", zap.Error(err))
	}
	return tlsCert, err
}

func (p *Proxy) httpsHandle(w http.ResponseWriter, req *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		p.log.Error("Hijack failed")
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hj.Hijack()
	if err != nil {
		p.log.Error("Hijack failed")
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	_, err = clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		p.log.Error("Failed to send 200 OK", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		p.log.Error("Failed splitting host/port", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	tlsCert, err := p.genCert(host)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		Certificates:     []tls.Certificate{tlsCert},
	}
	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	err = tlsConn.Handshake()
	if err != nil {
		p.log.Error("Handshake failed", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	connReader := bufio.NewReader(tlsConn)
	r, err := http.ReadRequest(connReader)
	if err == io.EOF {
		return
	} else if err != nil {
		p.log.Error("Failed to read request from client", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	reqBody, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(reqBody))

	changeRequestToTarget(r, req.Host)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		p.log.Error("Failed to send request to host", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(reqBody))
	err = p.saveRequest(resp, r)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	if err := resp.Write(tlsConn); err != nil {
		p.log.Error("Failed to send response to client", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	p.log.Debug("SUCCESS", zap.String("remote_addr", req.RemoteAddr), zap.Int("code", resp.StatusCode))
}

func (p *Proxy) saveRequest(resp *http.Response, r *http.Request) error {
	reqBody, _ := io.ReadAll(r.Body)
	r.Body.Close()

	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		r.Body = io.NopCloser(bytes.NewReader(reqBody))
		err := r.ParseForm()
		if err != nil {
			p.log.Error("Parse form failed", zap.Error(err))
			return err
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		p.log.Error("Read raw body failed", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(respBody))
		if err != nil {
			p.log.Error("Create gzip reader failed", zap.Error(err))
			return err
		}
		defer reader.Close()

		respBody, err = io.ReadAll(reader)
		if err != nil {
			p.log.Error("Read uncompressed data", zap.Error(err))
			return err
		}
	}

	cookies := parseCookies(r.Cookies())
	r.Header.Del("Cookie")
	return p.rep.Create(&requests.Request{
		Req: requests.Req{
			Method:     r.Method,
			Path:       r.URL.String(),
			GetParams:  r.URL.Query(),
			Headers:    r.Header,
			Cookies:    cookies,
			PostParams: r.PostForm,
			Body:       string(reqBody),
		},
		Resp: requests.Resp{
			Code:    resp.StatusCode,
			Message: resp.Status,
			Headers: resp.Header,
			Body:    string(respBody),
		},
	})
}
