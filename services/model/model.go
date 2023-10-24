package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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
