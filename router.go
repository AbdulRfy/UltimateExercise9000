package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"ultimate.com/exercise/authguard"
)

func setupRouter() {

	routerListenAddress := os.Getenv("RouterListenAddress")

	router := mux.NewRouter()

	router.HandleFunc("/register", registerUser).Methods("POST")
	router.HandleFunc("/login", loginUser).Methods("GET")
	router.Handle("/task", authguard.IsAuthorized(addTask)).Methods("POST")
	router.Handle("/task/{id:[0-9]+}", authguard.IsAuthorized(editTask)).Methods("PUT")
	router.Handle("/task/{id:[0-9]+}", authguard.IsAuthorized(deleteTask)).Methods("DELETE")
	router.Handle("/task", authguard.IsAuthorized(getAllUserTasks)).Methods("GET")
	router.Handle("/task/assign", authguard.IsAuthorized(assignTask)).Methods("PUT")

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(routerListenAddress, handler))
}
