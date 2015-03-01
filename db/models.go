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

type Tree interface {
	GetChildrens(*DB) error
	Childrens() []*Comment
}

type Item struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Brhub        Brhub         `json:"brhub"`
	Author       string        `json:"author"`
	CommentCount int           `json:"commentCount"`
	Comments     []*Comment    `json:"comments" bson:"-"`
	Date         int64         `json:"date"`
	Up           int           `json:"up"`
	Down         int           `json:"down"`
	Starred      bool          `json:"starred"`
}

func (item *Item) GetChildrens(db *DB) error {
	return db.C("comments").Find(bson.M{"item": item.Id, "parent": nil}).Sort("-date").All(&item.Comments)
}

func (item *Item) Childrens() []*Comment {
	return item.Comments
}

func NewItem() *Item {
	return &Item{
		Id:   bson.NewObjectId(),
		Date: time.Now().UnixNano(),
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
	Date     int64         `json:"date"`
	Up       int           `json:"up"`
	Down     int           `json:"down"`
	Parent   bson.ObjectId `json:"parent" bson:",omitempty"`
	Comments []*Comment    `json:"comments" bson:"-"`
	Item     bson.ObjectId `json:"-"`
}

func (com *Comment) GetChildrens(db *DB) error {
	return db.C("comments").Find(bson.M{"parent": com.Id}).Sort("-date").All(&com.Comments)
}

func (com *Comment) Childrens() []*Comment {
	return com.Comments
}

func NewComment() *Comment {
	return &Comment{
		Id:   bson.NewObjectId(),
		Date: time.Now().UnixNano(),
	}
}
