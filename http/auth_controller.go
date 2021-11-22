package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mattmeyers/heimdall/auth"
)

type AuthController struct {
	Service auth.Service
}

func (c *AuthController) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/auth", c.handleAuth)
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
