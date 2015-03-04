package handlers

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/klacabane/brhub/db"
	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
)

var (
	ctx       *Context
	user_test *db.User
)

func TestMain(m *testing.M) {
	session := db.MainSession()
	defer session.Close()

	var err error
	if user_test, err = session.DB().NewUser("foo", "bar"); err != nil {
		log.Fatal(err)
	}

	ctx = &Context{Session: session}

	m.Run()
}

func TestAuth(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8000/auth", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Form = url.Values{}
	req.Form.Add("name", "foo")
	req.Form.Add("password", "baz")

	code, data, err := Auth(ctx, web.C{}, req)
	assert.Equal(t, err, db.ErrFailAuth)
	assert.Equal(t, code, 401)
	assert.Nil(t, data)

	req.Form.Set("password", "bar")
	code, data, err = Auth(ctx, web.C{}, req)
	assert.Nil(t, err)
	assert.Equal(t, code, 200)
	assert.NotNil(t, data)

	user, ok := data.(*db.User)
	assert.True(t, ok)
	assert.True(t, len(user.Token) > 0)
}

func TestTimeline(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8000/api/timeline/0/5", nil)
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
	assert.Equal(t, code, 200)
	assert.Equal(t, data.([]*db.Item), []*db.Item{})
}
