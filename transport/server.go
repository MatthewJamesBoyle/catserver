package transport

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router(handler HttpHandler) *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/", handler.Get).Methods(http.MethodGet)
	return m
}
