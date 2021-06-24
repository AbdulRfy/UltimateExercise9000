package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func setupRouter() {

	routerListenAddress := os.Getenv("RouterListenAddress")

	router := mux.NewRouter()

	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("GET")
	router.Handle("/addTask", isAuthorized(addTask)).Methods("POST")
	router.Handle("/editTask/{id:[0-9]+}", isAuthorized(editTask)).Methods("PUT")
	router.Handle("/deleteTask/{id:[0-9]+}", isAuthorized(deleteTask)).Methods("DELETE")
	router.Handle("/getAllUserTasks/{userId:[0-9]+}", isAuthorized(getAllUserTasks)).Methods("GET")
	router.Handle("/assignTask", isAuthorized(assignTask)).Methods("PUT")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(routerListenAddress, handler))
}
