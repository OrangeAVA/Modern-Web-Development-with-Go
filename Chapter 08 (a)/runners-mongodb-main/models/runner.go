package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Runner struct {
	ID           string             `json:"id" bson:"-"`
	ObjectID     primitive.ObjectID `json:"-" bson:"_id"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Age          int                `json:"age,omitempty" bson:"age,omitempty"`
	IsActive     bool               `json:"is_active"`
	Country      string             `json:"country"`
	PersonalBest string             `json:"personal_best,omitempty" bson:"personalbest,omitempty"`
	SeasonBest   string             `json:"season_best,omitempty" bson:"seasonbest,omitempty"`
	Results      []*Result          `json:"results,omitempty" bson:"results,omitempty"`
}
