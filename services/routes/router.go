package routes

import (
	"cricket/services/controller"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/user", controller.CreatePlayer()).Methods("POST")      //add this
	router.HandleFunc("/user/{userId}", controller.GetPlayer()).Methods("GET") //add this
	router.HandleFunc("/user/{userId}", controller.UpdatePlayer()).Methods("PUT")
	router.HandleFunc("/users", controller.GetPlayers()).Methods("GET")
}
