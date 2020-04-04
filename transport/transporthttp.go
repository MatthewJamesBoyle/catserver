package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/matthewjamesboyle/catserver/internal/cat"
	"log"
	"net/http"
)

type HttpHandler struct {
	c cat.Servicer
}

func NewHttpHandler(c cat.Servicer) (*HttpHandler, error) {
	if c == nil {
		return nil, errors.New("nil servicer")
	}
	return &HttpHandler{c: c}, nil
}

func (h HttpHandler) Get(w http.ResponseWriter, req *http.Request) {
	c, err := h.c.GetImageAndFact(req.Context())
	if err != nil {
		log.Println(fmt.Errorf("error: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
	return
}
