package db

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"code.google.com/p/go.crypto/bcrypt"

	"labix.org/v2/mgo/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Password []byte        `json:"-"`
	Token    string        `json:"token"`
}

func (user *User) SetPassword(pw string) (err error) {
	user.Password, err = bcrypt.GenerateFromPassword([]byte(pw), 12)
	return
}

func (user *User) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
}

func (user *User) GenerateToken() (err error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["id"] = user.Id
	token.Claims["email"] = user.Email
	token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	user.Token, err = token.SignedString("secretkey")
	return
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
	err := db.C("comments").Find(bson.M{"item": item.Id, "parent": nil}).Sort("-date").All(&item.Comments)
	if item.Comments == nil {
		item.Comments = []*Comment{}
	}
	return err
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
	err := db.C("comments").Find(bson.M{"parent": com.Id}).Sort("-date").All(&com.Comments)
	if com.Comments == nil {
		com.Comments = []*Comment{}
	}
	return err
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
