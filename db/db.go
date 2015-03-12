package db

import (
	"fmt"
	"log"
	"sort"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	ErrFailAuth     = fmt.Errorf("auth failed")
	ErrUserNotFound = fmt.Errorf("user not found")
)

type Update map[string]interface{}

type Session struct {
	S *mgo.Session
}

func MainSession(addr string) *Session {
	session, err := mgo.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	ensureIndex := func(col string, key ...string) {
		if err = session.DB("brhub").C(col).EnsureIndexKey(key...); err != nil {
			log.Fatal(err)
		}
	}

	ensureIndex("brhubs", "name")
	ensureIndex("items", "-date")
	ensureIndex("comments", "-date", "parent", "item")
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

func (db *DB) AuthWithToken(name string, pw string) (*User, error) {
	var user *User
	err := db.C("users").Find(bson.M{"name": name}).One(&user)
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

func (db *DB) NewUser(name, password string) (*User, error) {
	user := &User{Id: bson.NewObjectId(), Name: name}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	err = user.GenerateToken()
	if err != nil {
		return nil, err
	}
	err = db.C("users").Insert(user)

	return user, err
}

func (db *DB) Star(userId, itemId bson.ObjectId) (Update, error) {
	err := db.C("users").UpdateId(userId, bson.M{"$addToSet": bson.M{"stars": itemId}})
	return Update{"starred": true}, err
}

func (db *DB) Unstar(userId, itemId bson.ObjectId) (Update, error) {
	err := db.C("users").UpdateId(userId, bson.M{"$pull": bson.M{"stars": itemId}})
	return Update{"starred": false}, err
}

func (db *DB) Timeline(userId bson.ObjectId, skip, limit int) ([]*Item, bool, error) {
	user, err := db.User(userId)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, false, ErrUserNotFound
		}
		return nil, false, err
	}

	items := make([]*Item, 0)
	err = db.C("items").Find(nil).Skip(skip).Limit(limit + 1).Sort("-date").All(&items)
	if err != nil {
		return nil, false, err
	}

	if len(user.Stars) > 0 {
		for _, item := range items {
			index := sort.Search(len(user.Stars), func(i int) bool {
				return user.Stars[i] >= item.Id
			})

			item.Starred = index < len(user.Stars) &&
				user.Stars[index] == item.Id
		}
	}

	var hasmore bool
	if len(items) > limit {
		items = items[:limit]
		hasmore = true
	}
	return items, hasmore, err
}

func (db *DB) AddBrhub(b *Brhub) error {
	return db.C("brhubs").Insert(b)
}

func (db *DB) AllBrhubs() ([]Brhub, error) {
	var all []Brhub
	err := db.C("brhubs").Find(nil).All(&all)
	return all, err
}

func (db *DB) Brhub(name string) (Brhub, error) {
	var b Brhub
	err := db.C("brhubs").Find(bson.M{"name": name}).One(&b)
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

func (db *DB) Items(brhubName string, skip, limit int) ([]*Item, bool, error) {
	items := make([]*Item, 0)
	err := db.C("items").
		Find(bson.M{"brhub": brhubName}).
		Skip(skip).Limit(limit + 1).
		Sort("-date").
		All(&items)
	if err != nil {
		return nil, false, err
	}

	var hasmore bool
	if len(items) > limit {
		items = items[:limit]
		hasmore = true
	}
	return items, hasmore, nil
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

func (db *DB) Upvote(itemId bson.ObjectId) (Update, error) {
	var item Item
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"upvote": 1}},
		Upsert:    false,
		ReturnNew: true,
	}
	_, err := db.C("items").FindId(itemId).Apply(change, &item)

	return Update{"upvote": item.Upvote}, err
}

func (db *DB) Comment(id bson.ObjectId) (Comment, error) {
	var c Comment
	err := db.C("comments").FindId(id).One(&c)
	return c, err
}

func (db *DB) AddComment(c *Comment) error {
	return db.C("comments").Insert(c)
}
