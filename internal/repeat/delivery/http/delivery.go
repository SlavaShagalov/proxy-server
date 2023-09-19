package http

import (
	"github.com/SlavaShagalov/proxy-server/internal/pkg/constants"
	"github.com/SlavaShagalov/proxy-server/internal/requests"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type delivery struct {
	rep    requests.Repository
	client *http.Client
	log    *zap.Logger
}

func RegisterHandlers(mux *mux.Router, rep requests.Repository, log *zap.Logger) {
	proxyUrl, _ := url.Parse(constants.ProxyURL)

	del := delivery{
		rep: rep,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		},
		log: log,
	}

	mux.HandleFunc("/repeat/{id}", del.repeat).Methods(http.MethodGet)
}

func (del *delivery) repeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, err := del.rep.Get(id)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	del.log.Debug("REPEAT", zap.String("id", req.ID))

	url, _ := url.Parse(req.Req.Path)

	_, err = del.client.Do(&http.Request{
		Method: req.Req.Method,
		URL:    url,
		Header: req.Req.Headers,
		Body:   io.NopCloser(strings.NewReader(req.Req.Body)),
	})
	if err != nil {
		del.log.Error("Repeat request failed", zap.Error(err))
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	del.log.Debug("SUCCESS", zap.String("id", req.ID))
}
