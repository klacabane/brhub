package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"labix.org/v2/mgo"
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
	id := c.URLParams["id"]
	if etag := r.Header.Get("If-None-Match"); len(etag) > 0 {
	}

	if !bson.IsObjectIdHex(id) {
		return 400, nil, fmt.Errorf("invalid id")
	}
	itemId := bson.ObjectIdHex(id)

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	item, err := session.DB().Item(itemId)
	if err != nil {
		if err == mgo.ErrNotFound {
			status = 404
		} else {
			status = 500
		}
	}
	return status, item, err
}

func CreateItem(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(r.FormValue("brhub")) {
		return 400, nil, fmt.Errorf("invalid brhub")
	}

	session := appCtx.SessionClone()
	defer session.Close()

	brhub, err := session.DB().Brhub(bson.ObjectIdHex(r.FormValue("brhub")))
	if err != nil {
		if err == mgo.ErrNotFound {
			return 400, nil, fmt.Errorf("invalid brhub")
		}
		return 500, nil, err
	}

	item := db.NewItem()
	item.Title = r.FormValue("title")
	item.Content = r.FormValue("content")
	item.Brhub = brhub

	status := 201
	err = session.DB().AddItem(item)
	if err != nil {
		status = 500
	}
	return status, item, err
}

func CreateBrhub(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	session := appCtx.SessionClone()
	defer session.Close()

	name := r.FormValue("name")
	if exists, err := session.DB().BrhubExists(name); err != nil {
		return 500, nil, err
	} else if exists {
		return 400, nil, fmt.Errorf("brhub %s already exists", name)
	}

	b := &db.Brhub{
		Id:   bson.NewObjectId(),
		Name: name,
	}

	status := 201
	err := session.DB().AddBrhub(b)
	if err != nil {
		status = 500
	}
	return status, b, err
}

func CreateComment(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	var itemId, parentId bson.ObjectId
	if !bson.IsObjectIdHex(r.FormValue("item")) {
		return 400, nil, fmt.Errorf("invalid item")
	}
	itemId = bson.ObjectIdHex(r.FormValue("item"))

	if parentstr := r.FormValue("parent"); len(parentstr) > 0 {
		if !bson.IsObjectIdHex(parentstr) {
			return 400, nil, fmt.Errorf("invalid parent")
		}
		parentId = bson.ObjectIdHex(parentstr)
	}

	session := appCtx.SessionClone()
	defer session.Close()

	if len(parentId) > 0 {
		_, err := session.DB().Comment(parentId)
		if err != nil {
			if err == mgo.ErrNotFound {
				return 400, nil, fmt.Errorf("invalid parent")
			}
			return 500, nil, fmt.Errorf(http.StatusText(500))
		}

	}

	_, err := session.DB().Item(itemId)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 400, nil, fmt.Errorf("invalid item")
		}
		return 500, nil, fmt.Errorf(http.StatusText(500))
	}

	comment := db.NewComment()
	comment.Parent = parentId
	comment.Item = itemId
	comment.Content = r.FormValue("content")

	status := 201
	err = session.DB().AddComment(comment)
	if err != nil {
		status = 500
	}
	return status, comment, err
}