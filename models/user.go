package models

import "labix.org/v2/mgo/bson"

type User struct {
	Id       bson.ObjectId `json: "id" bson:"_id,omitempty"`
	Name     string
	Mail     string
	Password []byte
}
