package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mattmeyers/heimdall/user"
)

type UserController struct {
	Service user.Service
}

func (c *UserController) Register(router *httprouter.Router) {
	router.HandlerFunc("POST", "/users", c.RegisterUser)
	router.HandlerFunc("GET", "/users/:id", c.GetByID)
	router.HandlerFunc("POST", "/login", c.Login)
}

type registrationBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body registrationBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, err := c.Service.Register(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf(`{"id":%d}`, id)))
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var body loginBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = c.Service.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(200)
}

func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	u, err := c.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	body, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)
}
