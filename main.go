package main 

import (
	"fmt"
	"github.com/gorilla/mux"
	"rest_full_api/controller/authcontroller"
	"rest_full_api/models"
	"rest_full_api/controller/productcontroller"
	"rest_full_api/middlewares"
	"net/http"
	"log"
)

func main() {
	models.ConnectDatabase()
	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.LogOut).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/product", productcontroller.Index).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	fmt.Println("app running in port : http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}