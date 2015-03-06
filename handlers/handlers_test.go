package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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

func TestFailAuth(t *testing.T) {
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

func TestSuccessAuth(t *testing.T) {
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
	assert.Equal(t, []*db.Item{}, data.([]*db.Item))
}

func TestCreateBrhub(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/b/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("name", "test_brhub")

	code, data, err := CreateBrhub(ctx, web.C{}, req)
	assert.Nil(t, err)
	assert.Equal(t, 201, code)

	brhub_test = data.(*db.Brhub)
	assert.Equal(t, "test_brhub", brhub_test.Name)
}

func TestCreateItem(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/items/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("brhub", "foo")
	req.Form.Add("title", "foobar")
	req.Form.Add("content", "lorem ipsum")

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

	req.Form.Set("brhub", brhub_test.Id.Hex())
	code, data, err = CreateItem(ctx, c, req)
	assert.Equal(t, 201, code)
	assert.Nil(t, err)

	item_test = data.(*db.Item)
	assert.Equal(t, "foobar", item_test.Title)
	assert.Equal(t, []*db.Comment{}, item_test.Comments)
}

func TestCreateComment(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/items/:id/comments", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("item", brhub_test.Id.Hex())
	req.Form.Add("content", "lorem ipsum")

	c := web.C{
		Env: map[string]interface{}{
			"user": db.Author{Id: user_test.Id, Name: user_test.Name},
		},
	}

	code, data, err := CreateComment(ctx, c, req)
	assert.Equal(t, 422, code)
	assert.NotNil(t, err)

	req.Form.Set("item", item_test.Id.Hex())
	code, data, err = CreateComment(ctx, c, req)
	assert.Equal(t, 200, code)
	assert.Nil(t, err)

	comment := data.(*db.Comment)
	assert.Equal(t, "lorem ipsum", comment.Content)
	item_test.Comments = append(item_test.Comments, comment)

	req.Form.Add("parent", comment.Id.Hex())
	code, data, err = CreateComment(ctx, c, req)
	assert.Equal(t, 200, code)
	assert.Nil(t, err)

	comment.Comments = append(comment.Comments, data.(*db.Comment))
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
