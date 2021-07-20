package router

import (
	"crud/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api", middleware.Hello).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return r
}
