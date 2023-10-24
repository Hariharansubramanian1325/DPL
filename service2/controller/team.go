package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cricket/service2/configs"
	"cricket/service2/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teamsCollection *mongo.Collection = configs.GetCollection(configs.DB, "teams")
var playerCollection *mongo.Collection = configs.GetCollection(configs.DB, "players")

func GetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var teams []model.Team

		cur, err := teamsCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var team model.Team
			err := cur.Decode(&team)
			if err != nil {
				log.Fatal(err)
			}

			teams = append(teams, team)

		}
		log.Println("reached!!!!!")
		json.NewEncoder(w).Encode(teams)
	}
}

func GetTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		teamID := params["teamId"]

		// if err != nil {
		// 	http.Error(w, "ID not found", http.StatusNotFound)
		// 	return
		// }
		var team model.Team
		objId, _ := primitive.ObjectIDFromHex(teamID)
		err1 := teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&team)
		if err1 != nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(team)
	}
}

func CreateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var team model.Team
		team.ID = primitive.NewObjectID()
		err := json.NewDecoder(r.Body).Decode(&team)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := teamsCollection.InsertOne(context.TODO(), team)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func UpdateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		teamID := params["teamId"]
		objId, _ := primitive.ObjectIDFromHex(teamID)
		var team model.Team

		err := json.NewDecoder(r.Body).Decode(&team)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"ID": objId}
		update := bson.M{
			"$set": bson.M{
				"ID":      objId,
				"name":    team.Name,
				"captain": team.Captain,
			},
		}

		_, err = teamsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "team not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(team)
	}
}
func Addplayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		log.Println("reached")
		playerId := params["playerid"]
		teamId := params["teamid"]
		objId, _ := primitive.ObjectIDFromHex(teamId)
		// objId2,_ := primitive.ObjectIDFromHex(playerId)
		var team model.Team
		err1 := teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&team)
		if err1 != nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}
		for _, player := range team.Players {
			if player == playerId {
				http.Error(w, "Player already exists", http.StatusNotFound)
				return
			}
		}
		team.Players = append(team.Players, playerId)
		filter := bson.M{"ID": objId}
		update := bson.M{
			"$set": bson.M{
				"ID":      objId,
				"name":    team.Name,
				"captain": team.Captain,
				"players": team.Players,
			},
		}
		_, err := teamsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "team not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(team)
	}
}
func Removeplayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		log.Println("reached")
		playerId := params["playerid"]
		teamId := params["teamid"]
		objId, _ := primitive.ObjectIDFromHex(teamId)
		var team model.Team
		err1 := teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&team)
		if err1 != nil {
			http.Error(w, "team not found", http.StatusNotFound)
			return
		}
		var c int = 0
		for _, player := range team.Players {
			log.Println(player)
			c++
		}
		var player2 []string
		var i int = 0
		for _, player := range team.Players {
			if player == playerId {
				continue
			} else {
				player2 = append(player2, player)
				i++
			}
		}
		filter := bson.M{"ID": objId}
		update := bson.M{
			"$set": bson.M{
				"ID":      objId,
				"name":    team.Name,
				"captain": team.Captain,
				"players": player2,
			},
		}
		if i == c {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}
		up, err := teamsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "team not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(up)
	}
}
func GetPlayers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		teamID := params["teamid"]
		objID, _ := primitive.ObjectIDFromHex(teamID)
		var team model.Team
		err := teamsCollection.FindOne(context.TODO(), bson.M{"ID": objID}).Decode(&team)
		if err != nil {
			http.Error(w, "Team not found", http.StatusNotFound)
			return
		}
		type PlayerResponse struct {
			Captain model.Player   `json:"captain"`
			Players []model.Player `json:"players"`
		}

		var captainPlayer model.Player
		captainID, _ := primitive.ObjectIDFromHex(team.Captain)
		err = playerCollection.FindOne(context.TODO(), bson.M{"ID": captainID}).Decode(&captainPlayer)
		if err != nil {
			http.Error(w, "Captain not found", http.StatusNotFound)
			return
		}
		response := PlayerResponse{
			Captain: captainPlayer,
			Players: []model.Player{},
		}

		for _, playerID := range team.Players {
			if playerID != team.Captain {
				var player model.Player
				playerID, _ := primitive.ObjectIDFromHex(playerID)
				err := playerCollection.FindOne(context.TODO(), bson.M{"ID": playerID}).Decode(&player)
				if err != nil {
					http.Error(w, "Player not found", http.StatusNotFound)
					return
				}
				response.Players = append(response.Players, player)
			}
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
