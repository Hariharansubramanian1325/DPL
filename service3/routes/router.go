package routes

import (
	"cricket/service3/controller"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/generate", controller.GenerateMatches()).Methods("POST")
	router.HandleFunc("/play", controller.PlayMatch()).Methods("POST")
	router.HandleFunc("/match/{matchid}", controller.Matchdetail()).Methods("GET")
	router.HandleFunc("/final", controller.Final()).Methods("GET")
	router.HandleFunc("/match", controller.Getallmatch()).Methods("GET")
	router.HandleFunc("/point", controller.Pointtable()).Methods("GET")
	router.HandleFunc("/final", controller.Finalinfo()).Methods("POST")
}
