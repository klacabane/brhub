package db

import (
	"github.com/klacabane/brhub/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Session struct {
	S *mgo.Session
}

func (s *Session) DB() *DB {
	return &DB{s.S.DB("brhub")}
}

func (s *Session) Close() {
	s.S.Close()
}

type DB struct {
	*mgo.Database
}

func (db *DB) Timeline(user bson.ObjectId) (posts []models.Post, err error) {
	err = db.C("posts").Find(bson.M{}).Sort("-date").All(&posts)
	return
}

func (db *DB) Posts(id bson.ObjectId) (posts []models.Post, err error) {
	err = db.C("posts").Find(bson.M{"brhub._id": id}).Sort("-date").All(&posts)
	return
}
