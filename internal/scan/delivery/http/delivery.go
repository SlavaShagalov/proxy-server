package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	pErrors "github.com/SlavaShagalov/proxy-server/internal/pkg/errors"
	"github.com/SlavaShagalov/proxy-server/internal/requests"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	stdUrl "net/url"
	"strings"
)

type delivery struct {
	rep    requests.Repository
	client *http.Client
	log    *zap.Logger
}

func RegisterHandlers(mux *mux.Router, rep requests.Repository, log *zap.Logger) {
	del := delivery{
		rep:    rep,
		client: &http.Client{},
		log:    log,
	}

	mux.HandleFunc("/scan/{id}", del.scan).Methods(http.MethodGet)
}

type response struct {
	GetInsecure  []string `json:"get_insecure"`
	PostInsecure []string `json:"post_insecure"`
}

func (del *delivery) scan(w http.ResponseWriter, r *http.Request) {
	const xssTest = `vulnerable'"><img src onerror=alert()>`

	vars := mux.Vars(r)
	id := vars["id"]

	dbReq, err := del.rep.Get(id)
	if err != nil {
		if err == pErrors.ErrRequestNotFound {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}
		del.log.Error("Get request error", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	del.log.Debug("SCAN", zap.String("id", dbReq.ID))

	url, err := stdUrl.Parse(dbReq.Req.Path)
	if err != nil {
		del.log.Error("Parse path failed", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	q := url.Query()

	apiResp := response{
		GetInsecure:  []string{},
		PostInsecure: []string{},
	}

	// scan get params
	for key, _ := range dbReq.Req.GetParams {
		oldValue := url.Query().Get(key)

		q.Set(key, xssTest)
		url.RawQuery = q.Encode()

		del.log.Debug("Try " + url.String())
		resp, err := del.client.Do(&http.Request{
			Method: dbReq.Req.Method,
			URL:    url,
			Header: dbReq.Req.Headers,
		})
		if err != nil {
			del.log.Error("Scan request failed", zap.Error(err))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			del.log.Error("Read raw body failed", zap.Error(err))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewReader(respBody))
			if err != nil {
				del.log.Error("Create gzip reader failed", zap.Error(err))
				return
			}
			defer reader.Close()

			respBody, err = io.ReadAll(reader)
			if err != nil {
				del.log.Error("Read uncompressed data", zap.Error(err))
				return
			}
		}

		if bytes.Contains(respBody, []byte(xssTest)) {
			apiResp.GetInsecure = append(apiResp.GetInsecure, key)
		}

		q.Set(key, oldValue)
		url.RawQuery = q.Encode()
	}

	postValues := stdUrl.Values{}
	for key, value := range dbReq.Req.PostParams {
		postValues.Set(key, value[0])
	}

	// scan post params
	for key, _ := range dbReq.Req.PostParams {
		oldValue := postValues.Get(key)
		postValues.Set(key, xssTest)

		req := &http.Request{
			Method: dbReq.Req.Method,
			URL:    url,
			Header: dbReq.Req.Headers,
			Body:   io.NopCloser(strings.NewReader(postValues.Encode())),
		}

		del.log.Debug("Try " + postValues.Encode())
		resp, err := del.client.Do(req)
		if err != nil {
			del.log.Error("Scan request failed", zap.Error(err))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			del.log.Error("Read raw body failed", zap.Error(err))
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if bytes.Contains(respBody, []byte(xssTest)) {
			apiResp.PostInsecure = append(apiResp.PostInsecure, key)
		}

		postValues.Set(key, oldValue)
	}

	data, err := json.Marshal(apiResp)
	if err != nil {
		del.log.Error("Marshal error", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
