package main

import (
	"cricket/service3/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/service3").Subrouter()
	routes.UserRoute(subrouter)
	go func() {
		fmt.Println("Server started successfully on port 8084")
		http.ListenAndServe(":8084", router)
	}()
	select {}

}
