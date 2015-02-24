package db

import (
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

func (db *DB) Timeline(user bson.ObjectId, skip, limit int) ([]Item, error) {
	items := make([]Item, limit)

	err := db.C("items").Find(bson.M{}).Skip(skip).Limit(limit).Sort("-date").All(&items)
	return items, err
}

func (db *DB) Items(id bson.ObjectId) ([]Item, error) {
	items := make([]Item, 0)

	err := db.C("items").Find(bson.M{"brhub._id": id}).Sort("-date").All(&items)
	return items, err
}

func (db *DB) CreateItem(item *Item) (err error) {
	err = db.C("items").Insert(item)
	return
}

func (db *DB) Item(id bson.ObjectId) (item Item, err error) {
	err = db.C("items").FindId(id).One(&item)
	return
}
