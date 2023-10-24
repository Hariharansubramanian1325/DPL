package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cricket/services/configs"
	"cricket/services/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var playersCollection *mongo.Collection = configs.GetCollection(configs.DB, "players")

func GetPlayers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var players []model.Player

		cur, err := playersCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var player model.Player
			err := cur.Decode(&player)
			if err != nil {
				log.Fatal(err)
			}

			players = append(players, player)

		}
		log.Println("reached!!!!!")
		json.NewEncoder(w).Encode(players)
	}
}

func GetPlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		playerID := params["userId"]

		// if err != nil {
		// 	http.Error(w, "ID not found", http.StatusNotFound)
		// 	return
		// }
		var player model.Player
		objId, _ := primitive.ObjectIDFromHex(playerID)
		err1 := playersCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&player)
		if err1 != nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(player)
	}
}

func CreatePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var player model.Player
		player.ID = primitive.NewObjectID()
		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := playersCollection.InsertOne(context.TODO(), player)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func UpdatePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		playerID := params["userId"]
		objId, _ := primitive.ObjectIDFromHex(playerID)
		var player model.Player

		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"ID": objId}
		update := bson.M{
			"$set": bson.M{
				"ID":         objId,
				"name":       player.Name,
				"jersey":     player.Jersey,
				"age":        player.Age,
				"prole":      player.Prole,
				"srole":      player.Srole,
				"batavg":     player.Batavg,
				"strikerate": player.Strikerate,
				"economy":    player.Economy,
				"matches":    player.Matches,
				"runs":       player.Runs,
				"wickets":    player.Wickets,
			},
		}

		_, err = playersCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, "Player not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(player)
	}
}
