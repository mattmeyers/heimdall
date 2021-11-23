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

type registerClientBody struct {
	RedirectURLs []string `json:"redirect_urls"`
}

func (c *ClientController) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var body registerClientBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	client, err := c.Service.Register(r.Context(), body.RedirectURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(client)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(resBody)
}
