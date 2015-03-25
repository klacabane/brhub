package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/klacabane/brhub/db"
	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
)

var (
	ctx        *Context
	user_test  *db.User
	brhub_test *db.Brhub
	item_test  *db.Item
)

func TestMain(m *testing.M) {
	session := db.MainSession("localhost")
	defer session.Close()
	defer session.DB().DropDatabase()

	var err error
	if user_test, err = session.DB().NewUser("foo", "bar"); err != nil {
		log.Fatal(err)
	}

	ctx = &Context{Session: session}

	m.Run()
}

func TestInvalidAuth(t *testing.T) {
	params := authp{Name: "foo", Password: "baz"}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	code, data, err := Auth(ctx, web.C{}, req)
	assert.Equal(t, err, db.ErrFailAuth)
	assert.Equal(t, 401, code)
	assert.Nil(t, data)
}

func TestValidAuth(t *testing.T) {
	params := authp{Name: "foo", Password: "bar"}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	code, data, err := Auth(ctx, web.C{}, req)
	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.NotNil(t, data)

	user, ok := data.(*db.User)
	assert.True(t, ok)
	assert.True(t, len(user.Token) > 0)
}

func TestTimeline(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/timeline/0/5", nil)
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
		URLParams: map[string]string{
			"skip":  "0",
			"limit": "5",
		},
	}

	code, data, err := Timeline(ctx, c, req)
	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.False(t, data.(page).Hasmore)
}

func TestCreateBrhub(t *testing.T) {
	params := authp{Name: "test_brhub"}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/b/", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	code, data, err := CreateBrhub(ctx, web.C{}, req)
	assert.Nil(t, err)
	assert.Equal(t, 201, code)

	brhub_test = data.(*db.Brhub)
	assert.Equal(t, "test_brhub", brhub_test.Name)
}

func TestInvalidCreateItem(t *testing.T) {
	params := itemp{"foo", "foobar", "lorem ipsum", db.TYPE_TEXT, ""}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/items/", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateItem(ctx, c, req)
	assert.Equal(t, 422, code)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid brhub", err.Error())
	assert.Nil(t, data)
}

func TestValidCreateItem(t *testing.T) {
	params := itemp{brhub_test.Name, "foobar", "lorem ipsum", db.TYPE_TEXT, ""}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/items/", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateItem(ctx, c, req)
	assert.Equal(t, 201, code)
	assert.Nil(t, err)

	item_test = data.(*db.Item)
	assert.Equal(t, "foobar", item_test.Title)
	assert.Equal(t, []*db.Comment{}, item_test.Comments)
}

func TestInvalidCreateComment(t *testing.T) {
	params := struct {
		Item, Content, Parent string
	}{"foo", "bar", ""}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/items/:id/comments", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateComment(ctx, c, req)
	assert.Equal(t, 422, code)
	assert.NotNil(t, err)
	assert.Nil(t, data)
}

func TestValidCreateComment(t *testing.T) {
	params := struct {
		Item, Content, Parent string
	}{item_test.Id.Hex(), "lorem ipsum", ""}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/comments/", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateComment(ctx, c, req)
	assert.Equal(t, 200, code)
	assert.Nil(t, err)

	comment := data.(*db.Comment)
	assert.Equal(t, "lorem ipsum", comment.Content)
	item_test.Comments = append(item_test.Comments, comment)
}

func TestCreateCommentWithParent(t *testing.T) {
	params := struct {
		Item, Content, Parent string
	}{item_test.Id.Hex(), "lorem ipsum", item_test.Comments[0].Id.Hex()}
	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/comments/", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateComment(ctx, c, req)
	assert.Equal(t, 200, code)
	assert.Nil(t, err)

	item_test.Comments[0].Comments = append(
		item_test.Comments[0].Comments, data.(*db.Comment))
	item_test.CommentCount = 2
}

func TestItem(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/items/:id", nil)
	if err != nil {
		log.Fatal(err)
	}

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
		URLParams: map[string]string{
			"id": item_test.Id.Hex(),
		},
	}

	code, data, err := Item(ctx, c, req)
	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.True(t, reflect.DeepEqual(data, item_test))
}
