package account

import (
	"fmt"

	authenticator "github.com/Timos-API/Authenticator"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	fmt.Println("Account routes registered")
	s := router.PathPrefix("/account").Subrouter()

	s.HandleFunc("/valid", authenticator.AuthMiddleware(nil, nil)).Methods("POST")

}
