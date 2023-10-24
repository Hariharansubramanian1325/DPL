package main

import (
	"cricket/service2/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/service2").Subrouter()
	routes.UserRoute(subrouter)
	go func() {
		fmt.Println("Server started successfully on port 8082")
		http.ListenAndServe(":8082", router)
	}()
	select {}
}
