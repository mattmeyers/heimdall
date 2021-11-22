package http

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mattmeyers/heimdall/auth"
)

type AuthController struct {
	Service auth.Service
}

func (c *AuthController) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/auth", c.handleAuth)
	router.HandlerFunc(http.MethodPost, "/login", c.handleLogin)
}

func (c *AuthController) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	token, err := c.Service.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	redirect, err := generateLoginRedirect(body.RedirectURL, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

type loginBody struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ClientID    string `json:"client_id"`
	RedirectURL string `json:"redirect_url"`
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
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		ClientID:    r.FormValue("client_id"),
		RedirectURL: r.FormValue("redirect_url"),
	}, nil
}

func generateLoginRedirect(redirectURL, token string) (string, error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "", err
	}

	params := u.Query()

	params.Set("token", token)
	params.Set("token_type", "JWT")
	params.Set("expires_in", strconv.Itoa(auth.JWTLifespan))

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func (c *AuthController) handleAuth(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("response_type") {
	case "token":
		clientID := r.URL.Query().Get("client_id")
		redirectURL := r.URL.Query().Get("redirect_url")

		tmpl, err := c.Service.ImplicitFlow(r.Context(), clientID, redirectURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		w.Write(tmpl)
	default:
		http.Error(w, "missing or invalid response_type", http.StatusBadRequest)
		return
	}
}
