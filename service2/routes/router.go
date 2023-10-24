package routes

import (
	"cricket/service2/controller"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/team", controller.CreateTeam()).Methods("POST")      //add this
	router.HandleFunc("/team/{teamId}", controller.GetTeam()).Methods("GET") //add this
	router.HandleFunc("/team/{teamId}", controller.UpdateTeam()).Methods("PUT")
	router.HandleFunc("/teams", controller.GetTeams()).Methods("GET")
	router.HandleFunc("/team/{teamid}/{playerid}", controller.Addplayer()).Methods("PUT")
	router.HandleFunc("/team/{teamid}/{playerid}", controller.Removeplayer()).Methods("DELETE")
	router.HandleFunc("/team/{teamid}", controller.GetPlayers()).Methods("POST")
}
