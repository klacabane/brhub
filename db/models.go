package db

import (
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"code.google.com/p/go.crypto/bcrypt"

	"labix.org/v2/mgo/bson"
)

var tagExp = regexp.MustCompile(`(?i)\[([a-z0-9$\s]+)\]`)

type Author struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string        `json:"name"`
}

type searchProps struct {
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

type User struct {
	Id       bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Password []byte          `json:"-"`
	Token    string          `json:"token" bson:"-"`
	Stars    []bson.ObjectId `json:"-"`
}

func (user *User) SetPassword(pw string) (err error) {
	user.Password, err = bcrypt.GenerateFromPassword([]byte(pw), 8)
	return
}

func (user *User) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
}

func (user *User) GenerateToken() (err error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["id"] = user.Id
	token.Claims["name"] = user.Name
	token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	user.Token, err = token.SignedString([]byte("secretkey"))
	return err
}

type Parent interface {
	GetChildrens(*DB) error
	Childrens() []*Comment
}

const (
	TYPE_TEXT = "text"
	TYPE_LINK = "link"
)

type Item struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title"`
	Link         string        `json:"link,omitempty" bson:"link,omitempty"`
	Content      string        `json:"content,omitempty" bson:"content,omitempty"`
	Tags         []string      `json:"tags"`
	Theme        Theme         `json:"theme"`
	Author       Author        `json:"author"`
	Comments     []*Comment    `json:"comments" bson:"-"`
	CommentCount int           `json:"commentCount" bson:"commentCount"`
	Date         int64         `json:"date"`
	Upvote       int           `json:"upvote"`
	Starred      bool          `json:"starred" bson:"-"`
	Type         string        `json:"type"`
	searchProps
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

func (item *Item) SetTitleAndTags(s string) {
	tags := tagExp.FindAllStringSubmatch(s, -1)
	item.Tags = make([]string, len(tags))
	for i, tag := range tags {
		item.Tags[i] = strings.TrimSpace(tag[1])
	}

	item.Title = strings.TrimSpace(
		tagExp.ReplaceAllString(s, ""))
}

func NewItem() *Item {
	return &Item{
		Id:       bson.NewObjectId(),
		Date:     time.Now().UnixNano(),
		Comments: []*Comment{},
		Tags:     []string{},
	}
}

type Theme struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	ColorHex string        `json:"color"`
	searchProps
}

type Comment struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Author   Author        `json:"author"`
	Content  string        `json:"content"`
	Date     int64         `json:"date"`
	Up       int           `json:"up"`
	Down     int           `json:"down"`
	Parent   bson.ObjectId `json:"parent" bson:",omitempty"`
	Comments []*Comment    `json:"comments" bson:"-"`
	Item     bson.ObjectId `json:"item"`
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
		Id:       bson.NewObjectId(),
		Date:     time.Now().UnixNano(),
		Comments: []*Comment{},
	}
}
