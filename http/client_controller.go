package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mattmeyers/heimdall/client"
)

type ClientController struct {
	Service client.Service
}

func (c *ClientController) Register(router *httprouter.Router) {
	router.HandlerFunc("GET", "/clients/:client_id", c.GetClientByID)
	router.HandlerFunc("POST", "/clients", c.RegisterClient)
}

func (c *ClientController) GetClientByID(w http.ResponseWriter, r *http.Request) {
	clientID := httprouter.ParamsFromContext(r.Context()).ByName("client_id")

	client, err := c.Service.Get(r.Context(), clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)
}

func (c *ClientController) RegisterClient(w http.ResponseWriter, r *http.Request) {
	client, err := c.Service.Register(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(body)
}
