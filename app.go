package main

import (
	"Timos-API/Accounts/auth"
	"Timos-API/Accounts/database"
	"Timos-API/Accounts/user"

	"log"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	auth.RegisterOAuth()
	database.Connect()

	router := mux.NewRouter()
	router.Use(routerMw)
	router.StrictSlash(true)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Origin"},
		AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	auth.RegisterRoutes(router)
	user.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	defer log.Fatal(server.ListenAndServe())

}

func routerMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/callback") {
			w.Header().Set("content-type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
