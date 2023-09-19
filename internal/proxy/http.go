package proxy

import (
	"go.uber.org/zap"
	"io"
	"net/http"
)

func copyHeaders(src *http.Response, dst http.ResponseWriter) {
	for key, values := range src.Header {
		for _, value := range values {
			dst.Header().Add(key, value)
		}
	}
}

func (p *Proxy) httpHandle(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	resp, err := p.client.Do(r)
	if err != nil {
		p.log.Error("Failed to send request to server", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	copyHeaders(resp, w)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		p.log.Error("Failed copying response body", zap.Error(err))
		http.Error(w, "Error copying response body", http.StatusInternalServerError)
		return
	}

	p.log.Debug("SUCCESS", zap.String("remote_addr", r.RemoteAddr), zap.Int("code", resp.StatusCode))
}
