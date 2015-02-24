package db

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type User struct {
	Id       bson.ObjectId `json: "id" bson:"_id,omitempty"`
	Name     string
	Mail     string
	Password []byte
}

type Item struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Brhub        Brhub         `json:"brhub"`
	Author       string        `json:"author"`
	CommentCount int           `json:"commentCount"`
	Comments     []Comment     `json:"comments"`
	Date         time.Time     `json:"date"`
}

func NewItem() *Item {
	return &Item{
		Id:       bson.NewObjectId(),
		Date:     time.Now(),
		Comments: make([]Comment, 0),
	}
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
