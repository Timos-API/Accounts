package main

import (
	"Timos-API/Accounts/persistence"
	"Timos-API/Accounts/service"
	"Timos-API/Accounts/transport"
	"context"
	"fmt"
	"os"

	"log"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	db := connectToDB()

	router := mux.NewRouter()
	router.Use(routerMw)
	router.StrictSlash(true)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Origin"},
		AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	// User module
	up := persistence.NewUserPersistor(db.Collection("user"))
	us := service.NewUserService(up)
	ut := transport.NewUserTransporter(us)
	ut.RegisterUserRoutes(router)

	// Auth module
	as := service.NewAuthService(us)
	at := transport.NewAuthTransporter(as)
	as.RegisterOAuth()
	at.RegisterAuthRoutes(router)

	server := &http.Server{
		Addr:         os.ExpandEnv("${host}:3000"),
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

func connectToDB() *mongo.Database {

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB")

	return c.Database("TimosAPI")
}
