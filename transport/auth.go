package transport

import (
	"Timos-API/Accounts/service"
	"fmt"
	"html/template"
	"net/http"

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

	s := router.PathPrefix("/auth").Subrouter()

	s.HandleFunc("/{provider}", gothic.BeginAuthHandler).Methods("GET")
	s.HandleFunc("/{provider}/callback", c.handleOAuthCallback).Methods("GET")

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
