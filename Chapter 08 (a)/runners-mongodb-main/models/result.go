package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Result struct {
	ID         string             `json:"id" bson:"id,omitempty"`
	ObjectID   primitive.ObjectID `json:"-" bson:"_id"`
	RunnerID   string             `json:"runner_id"`
	RaceResult string             `json:"race_result"`
	Location   string             `json:"location"`
	Position   int                `json:"position,omitempty"`
	Year       int                `json:"year"`
}
