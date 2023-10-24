package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Team struct {
	ID      primitive.ObjectID `bson:"ID,omitempty" validate:"required" json:"ID,omitempty"`
	Name    string             `bson:"name,omitempty" validate:"required"`
	Captain string             `bson:"captain,omitempty" validate:"required"`
	Players []string           `bson:"players,omitempty" validate:"required" json:"players,omitempty"`
}
type Player struct {
	ID         primitive.ObjectID `bson:"ID,omitempty" validate:"required" json:"ID,omitempty"`
	Name       string             `bson:"name,omitempty" validate:"required"`
	Jersey     string             `bson:"jersey,omitempty" validate:"required"`
	Age        int32              `bson:"age,omitempty" validate:"required"`
	Prole      string             `bson:"prole,omitempty" validate:"required"`
	Srole      string             `bson:"srole,omitempty" validate:"required"`
	Batavg     int32              `bson:"batavg,omitempty" validate:"required"`
	Strikerate float32            `bson:"strikerate,omitempty" validate:"required"`
	Economy    int32              `bson:"economy,omitempty" validate:"required"`
	Matches    int32              `bson:"matches,omitempty" validate:"required"`
	Runs       int32              `bson:"runs,omitempty" validate:"required"`
	Wickets    int32              `bson:"wickets,omitempty" validate:"required"`
}
type Matches struct {
	ID       primitive.ObjectID `bson:"ID,omitempty" validate:"required" json:"ID,omitempty"`
	Team1ID  string             `bson:"team1ID,omitempty" validate:"required"`
	Team2ID  string             `bson:"team2ID,omitempty" validate:"required"`
	WinnerID string             `bson:"winnerID,omitempty" validate:"required"`
	Venue    string             `bson:"venue,omitempty" validate:"required"`
	Date     string             `bson:"date,omitempty" validate:"required"`
	MoM      string             `bson:"mom,omitempty" validate:"required"`
	Runs     int32              `bson:"runs,omitempty" validate:"required" `
	Played   int32              `bson:"played,omitempty" validate:"required"`
}
