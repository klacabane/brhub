package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/klacabane/brhub/db"
	"github.com/zenazn/goji/web"
)

type authp struct {
	Name, Password string
}

type itemp struct {
	Brhub          bson.ObjectId
	Title, Content string
}

type commentp struct {
	Item    bson.ObjectId
	Content string
	// marshal as string to avoid error when empty
	Parent string

	parent bson.ObjectId
}

func Auth(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params authp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return 400, nil, fmt.Errorf(http.StatusText(400))
	}

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	user, err := session.DB().AuthWithToken(
		params.Name, params.Password)
	if err != nil {
		if err == db.ErrFailAuth {
			status = 401
		} else {
			status = 500
		}
	}
	return status, user, err
}

func Timeline(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
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
	userId := c.Env["user"].(db.Author).Id
	items, err := session.DB().Timeline(userId, int(skip), int(limit))
	if err != nil {
		status = 500
	}
	return status, items, err
}

func Items(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(c.URLParams["id"]) {
		return 422, nil, fmt.Errorf("invalid id")
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

func Item(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	id := c.URLParams["id"]
	if etag := r.Header.Get("If-None-Match"); len(etag) > 0 {
	}

	if !bson.IsObjectIdHex(id) {
		return 422, nil, fmt.Errorf("invalid id")
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

func CreateItem(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params itemp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return 422, nil, fmt.Errorf("invalid brhub")
	}

	session := appCtx.SessionClone()
	defer session.Close()

	brhub, err := session.DB().Brhub(params.Brhub)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 422, nil, fmt.Errorf("invalid brhub")
		}
		return 500, nil, err
	}

	item := db.NewItem()
	item.Title = params.Title
	item.Content = params.Content
	item.Brhub = brhub
	item.Author = c.Env["user"].(db.Author)

	status := 201
	err = session.DB().AddItem(item)
	if err != nil {
		status = 500
	}
	return status, item, err
}

func Upvote(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(c.URLParams["id"]) {
		return 422, nil, fmt.Errorf("invalid item")
	}

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	dif, err := session.DB().Upvote(bson.ObjectIdHex(c.URLParams["id"]))
	if err != nil {
		if err == mgo.ErrNotFound {
			status = 404
		} else {
			status = 500
		}
	}
	return status, dif, err
}

func Star(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(r.FormValue("item")) {
		return 422, nil, fmt.Errorf("invalid item")
	}

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	itemId := bson.ObjectIdHex(r.FormValue("item"))

	_, err := session.DB().Item(itemId)
	if err != nil {
		status = 500
		if err == mgo.ErrNotFound {
			status = 422
		}
		return status, nil, err
	}

	dif, err := session.DB().Star(
		c.Env["user"].(db.Author).Id,
		itemId)
	if err != nil {
		status = 500
	}
	return status, dif, err
}

func CreateBrhub(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params authp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return 400, nil, fmt.Errorf(http.StatusText(400))
	}

	session := appCtx.SessionClone()
	defer session.Close()

	if exists, err := session.DB().BrhubExists(params.Name); err != nil {
		return 500, nil, err
	} else if exists {
		return 422, nil, fmt.Errorf("brhub %s already exists", params.Name)
	}

	b := &db.Brhub{
		Id:   bson.NewObjectId(),
		Name: params.Name,
	}

	status := 201
	err := session.DB().AddBrhub(b)
	if err != nil {
		status = 500
	}
	return status, b, err
}

func CreateComment(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params commentp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return 422, nil, err
	}

	session := appCtx.SessionClone()
	defer session.Close()

	if len(params.Parent) > 0 {
		if !bson.IsObjectIdHex(params.Parent) {
			return 422, nil, fmt.Errorf("invalid parent")
		}
		params.parent = bson.ObjectIdHex(params.Parent)

		_, err := session.DB().Comment(params.parent)
		if err != nil {
			if err == mgo.ErrNotFound {
				return 422, nil, fmt.Errorf("invalid parent")
			}
			return 500, nil, fmt.Errorf(http.StatusText(500))
		}
	}

	_, err := session.DB().Item(params.Item)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 422, nil, fmt.Errorf("invalid item")
		}
		return 500, nil, fmt.Errorf(http.StatusText(500))
	}

	comment := db.NewComment()
	comment.Parent = params.parent
	comment.Item = params.Item
	comment.Content = params.Content
	comment.Author = c.Env["user"].(db.Author)

	status := 200
	err = session.DB().AddComment(comment)
	if err != nil {
		status = 500
	}
	return status, comment, err
}
