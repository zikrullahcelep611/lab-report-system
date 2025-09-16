package auth

import "github.com/gorilla/mux"

func RegisterAuthRoutes(routes *mux.Router, handler *AuthHandler){
	routes.HandleFunc("/login", handler.Login).Methods("POST")
	routes.HandleFunc("/logout", handler.Logout).Methods("POST")
}