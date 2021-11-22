package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	var err error

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		body, err = getloginBodyFromJSON(r)
	} else if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		body, err = getloginBodyFromForm(r)
	} else {
		http.Error(w, "Unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.Service.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(200)
}

func getloginBodyFromJSON(r *http.Request) (loginBody, error) {
	var body loginBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return loginBody{}, err
	}

	return body, nil
}

func getloginBodyFromForm(r *http.Request) (loginBody, error) {
	if err := r.ParseForm(); err != nil {
		return loginBody{}, err
	}

	return loginBody{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}, nil
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
