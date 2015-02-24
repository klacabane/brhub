package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"labix.org/v2/mgo/bson"

	"github.com/klacabane/brhub/db"
	"github.com/zenazn/goji/web"
)

func Timeline(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var skip, limit int64
	if limit, _ = strconv.ParseInt(c.URLParams["limit"], 10, 64); limit <= 0 || limit > 10 {
		limit = 10
	}
	if skip, _ = strconv.ParseInt(c.URLParams["skip"], 10, 64); skip < 0 {
		skip = 0
	}

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	items, err := session.DB().Timeline(bson.NewObjectId() /*userid*/, int(skip), int(limit))
	if err != nil {
		status = 500
	}
	return status, items, err
}

func Items(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(c.URLParams["id"]) {
		return 400, nil, fmt.Errorf("invalid id")
	}
	brhubId := bson.ObjectIdHex(c.URLParams["id"])

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	items, err := session.DB().Items(brhubId)
	if err != nil {
		status = 500
	}
	return status, items, err
}

func Item(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(c.URLParams["id"]) {
		return 400, nil, fmt.Errorf("invalid id")
	}
	itemId := bson.ObjectIdHex(c.URLParams["id"])

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	item, err := session.DB().Item(itemId)
	if err != nil {
		status = 500
	}
	return status, item, err
}

func CreateItem(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var item = db.NewItem()
	item.Title = r.FormValue("title")
	item.Content = r.FormValue("content")

	session := appCtx.SessionClone()
	defer session.Close()

	status := 201
	err := session.DB().CreateItem(item)
	if err != nil {
		status = 500
	}
	return status, item, err
}
