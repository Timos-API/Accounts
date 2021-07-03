package user

import (
	"fmt"

	authenticator "github.com/Timos-API/Authenticator"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	fmt.Println("User routes registered")
	s := router.PathPrefix("/user").Subrouter()

	s.HandleFunc("/valid", authenticator.AuthMiddleware(nil, nil)).Methods("POST")
	s.HandleFunc("/info/{id}", GetUserInfo).Methods("GET")
}
