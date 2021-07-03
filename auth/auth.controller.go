package auth

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func RegisterRoutes(router *mux.Router) {
	fmt.Println("Auth routes registered")
	s := router.PathPrefix("/auth").Subrouter()

	s.HandleFunc("/{provider}", gothic.BeginAuthHandler).Methods("GET")
	s.HandleFunc("/{provider}/callback", handleOAuthCallback).Methods("GET")
}
