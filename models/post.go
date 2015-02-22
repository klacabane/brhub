package models

import "labix.org/v2/mgo/bson"

type Post struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Brhub        Brhub         `json:"brhub"`
	CommentCount int           `json:"commentCount"`
	Comments     []Comment     `json:"comments"`
}

type Brhub struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string        `json:"name"`
}

type Comment struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author   string        `json:"author"`
	Content  string        `json:"content"`
	Comments []Comment     `json:"comments"`
}
