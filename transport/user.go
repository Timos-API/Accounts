package transport

import (
	"Timos-API/Accounts/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type UserTransporter struct {
	s *service.UserService
}

func NewUserTransporter(s *service.UserService) *UserTransporter {
	return &UserTransporter{s}
}

func (c *UserTransporter) RegisterUserRoutes(router *mux.Router) {

	router.HandleFunc("/user/info/{id}", c.getUserInfo).Methods("GET")

	fmt.Println("User routes registered")
}

func (c *UserTransporter) getUserInfo(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]

	if !ok {
		http.Error(w, "Missing param: id", http.StatusBadRequest)
		return
	}

	userInfo, err := c.s.GetUserInfo(req.Context(), id)

	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(userInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
