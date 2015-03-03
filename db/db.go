package db

import (
	"fmt"
	"log"
	"sort"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var ErrFailAuth = fmt.Errorf("auth failed")

type Session struct {
	S *mgo.Session
}

func MainSession() *Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}

	if err = session.DB("brhub").C("items").EnsureIndexKey("-date"); err != nil {
		log.Fatal(err)
	}
	if err = session.DB("brhub").C("comments").EnsureIndexKey("-date", "parent", "item"); err != nil {
		log.Fatal(err)
	}
	return &Session{S: session}
}

func (s *Session) DB() *DB {
	return &DB{s.S.DB("brhub")}
}

func (s *Session) Clone() *Session {
	return &Session{S: s.S.Clone()}
}

func (s *Session) Close() {
	s.S.Close()
}

type DB struct {
	*mgo.Database
}

func (db *DB) AuthWithToken(mail string, pw string) (*User, error) {
	var user *User
	err := db.C("users").Find(bson.M{"mail": mail}).One(&user)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrFailAuth
		}
		return nil, err
	}

	if err = user.ComparePassword(pw); err != nil {
		return nil, ErrFailAuth
	}
	err = user.GenerateToken()

	return user, err
}

func (db *DB) User(id bson.ObjectId) (*User, error) {
	var user *User
	err := db.C("users").FindId(id).One(&user)
	return user, err
}

func (db *DB) Timeline(userId bson.ObjectId, skip, limit int) ([]*Item, error) {
	user, err := db.User(userId)
	if err != nil {
		return nil, err
	}

	items := make([]*Item, limit)
	err = db.C("items").Find(nil).Skip(skip).Limit(limit).Sort("-date").All(&items)
	if err != nil {
		return nil, err
	}

	if len(user.Stars) > 0 {
		for _, item := range items {
			idstr := item.Id.String()
			index := sort.Search(len(user.Stars), func(i int) bool {
				return user.Stars[i].String() >= idstr
			})

			item.Starred = index < len(user.Stars) &&
				user.Stars[index].String() == idstr
		}
	}
	return items, err
}

func (db *DB) AddBrhub(b *Brhub) error {
	return db.C("brhubs").Insert(b)
}

func (db *DB) Brhub(id bson.ObjectId) (Brhub, error) {
	var b Brhub
	err := db.C("brhubs").Find(bson.M{"_id": id}).One(&b)
	return b, err
}

func (db *DB) BrhubExists(name string) (bool, error) {
	var b Brhub
	err := db.C("brhubs").Find(bson.M{"name": name}).One(&b)
	if err != nil && err == mgo.ErrNotFound {
		return false, nil
	}
	return true, err
}

func (db *DB) Items(brhubId bson.ObjectId) ([]Item, error) {
	items := make([]Item, 0)
	err := db.C("items").Find(bson.M{"brhub._id": brhubId}).Sort("-date").All(&items)
	return items, err
}

func (db *DB) AddItem(item *Item) error {
	return db.C("items").Insert(item)
}

func (db *DB) Item(id bson.ObjectId) (*Item, error) {
	var item *Item
	err := db.C("items").FindId(id).One(&item)
	if err != nil {
		return item, err
	}

	err = db.buildCommentTree(item)

	return item, err
}

func (db *DB) buildCommentTree(item *Item) error {
	errchan := make(chan error, 1)
	go db.commentTree(item, errchan)
	return <-errchan
}

func (db *DB) commentTree(parent Tree, c chan<- error) {
	err := parent.GetChildrens(db)
	if err != nil {
		c <- err
		return
	}

	clen := len(parent.Childrens())
	if clen == 0 {
		c <- nil
		return
	}

	errchan := make(chan error, clen)
	for _, child := range parent.Childrens() {
		go db.commentTree(child, errchan)
	}

	for i := 0; i < clen; i++ {
		errc := <-errchan
		if err == nil {
			err = errc
		}
	}
	c <- err
}

func (db *DB) Comment(id bson.ObjectId) (Comment, error) {
	var c Comment
	err := db.C("comments").FindId(id).One(&c)
	return c, err
}

func (db *DB) AddComment(c *Comment) error {
	return db.C("comments").Insert(c)
}
