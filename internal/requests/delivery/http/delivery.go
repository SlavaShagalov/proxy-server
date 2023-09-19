package http

import (
	"encoding/json"
	pRequests "github.com/SlavaShagalov/proxy-server/internal/requests"
	"github.com/gorilla/mux"
	"net/http"
)

type delivery struct {
	uc pRequests.Usecase
}

func RegisterHandlers(mux *mux.Router, uc pRequests.Usecase) {
	del := delivery{
		uc: uc,
	}

	mux.HandleFunc("/requests", del.list).Methods(http.MethodGet)
	mux.HandleFunc("/requests/{id}", del.get).Methods(http.MethodGet)
}

func (del *delivery) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, _ := del.uc.Get(id)

	data, _ := json.Marshal(req)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (del *delivery) list(w http.ResponseWriter, r *http.Request) {
	reqs, _ := del.uc.List()

	data, _ := json.Marshal(reqs)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
