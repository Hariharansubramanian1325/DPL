package main

import (
	"cricket/services/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/service1").Subrouter()
	routes.UserRoute(subrouter)
	go func() {
		fmt.Println("Server started successfully on port 8083")
		http.ListenAndServe(":8083", router)
	}()
	select {}

}
