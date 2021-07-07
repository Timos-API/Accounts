package transport

import (
	"Timos-API/Accounts/service"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Timos-API/authenticator"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type AuthTransporter struct {
	s *service.AuthService
}

func NewAuthTransporter(s *service.AuthService) *AuthTransporter {
	return &AuthTransporter{s}
}

func (c *AuthTransporter) RegisterAuthRoutes(router *mux.Router) {

	router.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", c.handleOAuthCallback).Methods("GET")
	router.HandleFunc("/auth/valid", authenticator.Middleware(nil, nil)).Methods("POST")

	fmt.Println("Auth routes registered")
}

func (c *AuthTransporter) handleOAuthCallback(w http.ResponseWriter, req *http.Request) {

	gothUser, err := gothic.CompleteUserAuth(w, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := c.s.UserSignedIn(req.Context(), gothUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, err := template.ParseFiles("auth.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp.Execute(w, token)

}
