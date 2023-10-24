package controller

import (
	"context"
	"cricket/service3/configs"
	"cricket/service3/model"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teamsCollection *mongo.Collection = configs.GetCollection(configs.DB, "teams")
var playerCollection *mongo.Collection = configs.GetCollection(configs.DB, "players")
var matchCollection *mongo.Collection = configs.GetCollection(configs.DB, "matches")

func GenerateMatches() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var teams []model.Team

		cur, err := teamsCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var team model.Team
			err := cur.Decode(&team)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			teams = append(teams, team)
		}
		matches := []model.Matches{}
		for i := 0; i < len(teams); i++ {
			for j := i + 1; j < len(teams); j++ {
				team1 := teams[i]
				team2 := teams[j]

				match := model.Matches{
					ID:       primitive.NewObjectID(),
					Team1ID:  team1.ID.Hex(),
					Team2ID:  team2.ID.Hex(),
					WinnerID: "TBD",
					Venue:    "TBD",
					Date:     "TBD",
					MoM:      "TBD",
					Runs:     -1,
					Played:   -1,
				}
				matches = append(matches, match)
			}
		}

		for _, match := range matches {
			_, err := matchCollection.InsertOne(context.TODO(), match)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		err = json.NewEncoder(w).Encode(matches)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PlayMatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var matches []model.Matches

		cur, err := matchCollection.Find(context.TODO(), bson.M{"played": -1})
		if err != nil {
			log.Printf("error1")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {
			var match model.Matches
			err := cur.Decode(&match)
			if err != nil {
				log.Printf("error2")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			matches = append(matches, match)
		}

		venues := []string{"chennai", "delhi", "mumbai", "kolkata"}

		rand.Seed(time.Now().UnixNano())

		for _, match := range matches {

			venueIndex := rand.Intn(len(venues))
			match.Venue = venues[venueIndex]

			winnerIndex := rand.Intn(2)
			if winnerIndex == 0 {
				match.WinnerID = match.Team1ID
			} else {
				match.WinnerID = match.Team2ID
			}
			objId, _ := primitive.ObjectIDFromHex(match.WinnerID)

			var winningTeam model.Team
			err := teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&winningTeam)
			if err != nil {
				log.Printf("error3")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			match.Runs = int32(rand.Intn(101) + 100)
			momIndex := rand.Intn(len(winningTeam.Players))
			momPlayerID := winningTeam.Players[momIndex]

			match.MoM = momPlayerID
			match.Date = time.Now().Format("2006-01-02")

			match.Played = 1

			_, err2 := matchCollection.UpdateOne(context.TODO(), bson.M{"ID": match.ID}, bson.M{"$set": match})
			if err2 != nil {
				log.Printf("error4")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = json.NewEncoder(w).Encode(matches)
		if err != nil {
			log.Printf("error5")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func Matchdetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		matchID := params["matchid"]
		objId, _ := primitive.ObjectIDFromHex(matchID)
		var match model.Matches
		err := matchCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		var team1 model.Team
		objId1, _ := primitive.ObjectIDFromHex(match.Team1ID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId1}).Decode(&team1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var team2 model.Team
		objId2, _ := primitive.ObjectIDFromHex(match.Team2ID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId2}).Decode(&team2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var winner model.Team
		objId4, _ := primitive.ObjectIDFromHex(match.WinnerID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId4}).Decode(&winner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var momPlayer model.Player
		objId3, _ := primitive.ObjectIDFromHex(match.MoM)
		err = playerCollection.FindOne(context.TODO(), bson.M{"ID": objId3}).Decode(&momPlayer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchDetails := map[string]interface{}{
			"MatchID":  match.ID.Hex(),
			"Date":     match.Date,
			"Venue":    match.Venue,
			"Runs":     match.Runs,
			"WinnerID": winner.Name,
			"Team1":    team1.Name,
			"Team2":    team2.Name,
			"MoM":      momPlayer.Name,
		}

		err = json.NewEncoder(w).Encode(matchDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func Getallmatch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cur, err := matchCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		for cur.Next(context.TODO()) {
			var match model.Matches
			err := cur.Decode(&match)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var team1 model.Team
			objId1, _ := primitive.ObjectIDFromHex(match.Team1ID)
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId1}).Decode(&team1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var team2 model.Team
			objId2, _ := primitive.ObjectIDFromHex(match.Team2ID)
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId2}).Decode(&team2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var winner model.Team
			objId4, _ := primitive.ObjectIDFromHex(match.WinnerID)
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId4}).Decode(&winner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var momPlayer model.Player
			objId3, _ := primitive.ObjectIDFromHex(match.MoM)
			err = playerCollection.FindOne(context.TODO(), bson.M{"ID": objId3}).Decode(&momPlayer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			matchDetails := map[string]interface{}{
				"MatchID":  match.ID.Hex(),
				"Date":     match.Date,
				"Venue":    match.Venue,
				"Runs":     match.Runs,
				"WinnerID": winner.Name,
				"Team1":    team1.Name,
				"Team2":    team2.Name,
				"MoM":      momPlayer.Name,
			}

			err = json.NewEncoder(w).Encode(matchDetails)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
	}
}

type PointTable struct {
	TeamName string
	Matches  int
	Points   int
}

func Pointtable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cur, err := matchCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		pointTable := make(map[string]*PointTable)

		for cur.Next(context.TODO()) {
			var match model.Matches
			err := cur.Decode(&match)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			objId, _ := primitive.ObjectIDFromHex(match.Team1ID)
			var team1 model.Team
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&team1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			objId2, _ := primitive.ObjectIDFromHex(match.Team2ID)
			var team2 model.Team
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId2}).Decode(&team2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, exists := pointTable[team1.Name]; !exists {
				pointTable[team1.Name] = &PointTable{TeamName: team1.Name}
			}
			if _, exists := pointTable[team2.Name]; !exists {
				pointTable[team2.Name] = &PointTable{TeamName: team2.Name}
			}
			if match.WinnerID == team1.ID.Hex() {
				pointTable[team1.Name].Matches++
				pointTable[team2.Name].Matches++
				pointTable[team1.Name].Points += 2
			} else if match.WinnerID == team2.ID.Hex() {
				pointTable[team2.Name].Matches++
				pointTable[team1.Name].Matches++
				pointTable[team2.Name].Points += 2
			}
		}
		pointTableList := []*PointTable{}
		for _, p := range pointTable {
			pointTableList = append(pointTableList, p)
		}

		err = json.NewEncoder(w).Encode(pointTableList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
func Final() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cur, err := matchCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		pointTable := make(map[string]*PointTable)

		for cur.Next(context.TODO()) {
			var match model.Matches
			err := cur.Decode(&match)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			objId, _ := primitive.ObjectIDFromHex(match.Team1ID)
			var team1 model.Team
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId}).Decode(&team1)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			objId2, _ := primitive.ObjectIDFromHex(match.Team2ID)
			var team2 model.Team
			err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId2}).Decode(&team2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, exists := pointTable[team1.Name]; !exists {
				pointTable[team1.Name] = &PointTable{TeamName: team1.Name}
			}
			if _, exists := pointTable[team2.Name]; !exists {
				pointTable[team2.Name] = &PointTable{TeamName: team2.Name}
			}
			if match.WinnerID == team1.ID.Hex() {
				pointTable[team1.Name].Matches++
				pointTable[team2.Name].Matches++
				pointTable[team1.Name].Points += 2
			} else if match.WinnerID == team2.ID.Hex() {
				pointTable[team2.Name].Matches++
				pointTable[team1.Name].Matches++
				pointTable[team2.Name].Points += 2
			}
		}
		pointTableList := []*PointTable{}
		for _, p := range pointTable {
			pointTableList = append(pointTableList, p)
		}

		var topTeams []*PointTable
		for _, p := range pointTableList {
			if len(topTeams) < 2 {
				topTeams = append(topTeams, p)
			} else if p.Points > topTeams[0].Points {
				topTeams[1] = topTeams[0]
				topTeams[0] = p
			} else if p.Points > topTeams[1].Points {
				topTeams[1] = p
			}
		}

		if len(topTeams) < 2 {
			http.Error(w, "Not enough teams for a match", http.StatusInternalServerError)
			return
		}
		var fteam1 model.Team
		var fteam2 model.Team
		err1 := teamsCollection.FindOne(context.TODO(), bson.M{"name": topTeams[0].TeamName}).Decode(&fteam1)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"name": topTeams[1].TeamName}).Decode(&fteam2)
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		match := model.Matches{
			ID:       primitive.NewObjectID(),
			Team1ID:  fteam1.ID.Hex(),
			Team2ID:  fteam2.ID.Hex(),
			WinnerID: "TBD",
			Venue:    "Dubai",
			Date:     time.Now().Format("2006-01-02"),
			MoM:      "TBD",
			Runs:     int32(rand.Intn(101) + 100),
			Played:   0,
		}

		if rand.Intn(2) == 0 {
			match.WinnerID = fteam1.ID.Hex()
			momIndex := rand.Intn(len(fteam1.Players))
			momPlayerID := fteam1.Players[momIndex]
			match.MoM = momPlayerID
		} else {
			momIndex := rand.Intn(len(fteam2.Players))
			momPlayerID := fteam2.Players[momIndex]
			match.MoM = momPlayerID
			match.WinnerID = fteam2.ID.Hex()
		}

		match.Played = 2

		_, err2 := matchCollection.InsertOne(context.TODO(), match)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func Finalinfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var match model.Matches
		err := matchCollection.FindOne(context.TODO(), bson.M{"played": 2}).Decode(&match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var team1 model.Team
		objId1, _ := primitive.ObjectIDFromHex(match.Team1ID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId1}).Decode(&team1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var team2 model.Team
		objId2, _ := primitive.ObjectIDFromHex(match.Team2ID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId2}).Decode(&team2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var winner model.Team
		objId4, _ := primitive.ObjectIDFromHex(match.WinnerID)
		err = teamsCollection.FindOne(context.TODO(), bson.M{"ID": objId4}).Decode(&winner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var momPlayer model.Player
		objId3, _ := primitive.ObjectIDFromHex(match.MoM)
		err = playerCollection.FindOne(context.TODO(), bson.M{"ID": objId3}).Decode(&momPlayer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchDetails := map[string]interface{}{
			"MatchID":  match.ID.Hex(),
			"Date":     match.Date,
			"Venue":    match.Venue,
			"Runs":     match.Runs,
			"WinnerID": winner.Name,
			"Team1":    team1.Name,
			"Team2":    team2.Name,
			"MoM":      momPlayer.Name,
		}

		err = json.NewEncoder(w).Encode(matchDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
