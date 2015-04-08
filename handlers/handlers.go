package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/klacabane/brhub/db"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zenazn/goji/web"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type authp struct {
	Name, Password string
}

type brhubp struct {
	Name, Color string
}

type itemp struct {
	Brhub, Title, Content, Type, Link string
}

type commentp struct {
	Item    bson.ObjectId
	Content string
	// marshal as string to avoid error when empty
	Parent string

	parent bson.ObjectId
}

type page struct {
	Hasmore bool       `json:"hasmore"`
	Items   []*db.Item `json:"items"`
}

func decodeParams(rc io.Reader, params interface{}) error {
	decoder := json.NewDecoder(rc)
	return decoder.Decode(&params)
}

func Auth(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params authp
	if err := decodeParams(r.Body, &params); err != nil {
		return 400, nil, errors.New(http.StatusText(400))
	}
	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	user, err := session.DB().AuthWithToken(params.Name, params.Password)
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

	userId := c.Env["user"].(db.Author).Id
	items, hasmore, err := session.DB().Timeline(userId, int(skip), int(limit))
	if err != nil {
		status := 500
		if err == db.ErrUserNotFound {
			status = 403
		}
		return status, nil, err
	}
	return 200, page{hasmore, items}, nil
}

func Items(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var skip, limit int64
	if limit, _ = strconv.ParseInt(c.URLParams["limit"], 10, 64); limit <= 0 || limit > 10 {
		limit = 10
	}
	if skip, _ = strconv.ParseInt(c.URLParams["skip"], 10, 64); skip < 0 {
		skip = 0
	}
	session := appCtx.SessionClone()
	defer session.Close()

	_, err := session.DB().Brhub(c.URLParams["name"])
	if err != nil {
		if err == mgo.ErrNotFound {
			return 422, nil, errors.New("invalid brhub")
		}
		return 500, nil, errors.New(http.StatusText(500))
	}
	items, hasmore, err := session.DB().Items(c.URLParams["name"], int(skip), int(limit))
	if err != nil {
		return 500, nil, err
	}
	return 200, page{hasmore, items}, nil
}

func Item(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	id := c.URLParams["id"]
	if etag := r.Header.Get("If-None-Match"); len(etag) > 0 {
	}

	if !bson.IsObjectIdHex(id) {
		return 422, nil, errors.New("invalid id")
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
		return 500, nil, errors.New(http.StatusText(500))
	}

	if params.Type != db.TYPE_TEXT && params.Type != db.TYPE_LINK {
		return 422, nil, errors.New("invalid type")
	}

	session := appCtx.SessionClone()
	defer session.Close()

	b, err := session.DB().Brhub(params.Brhub)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 422, nil, errors.New("invalid brhub")
		}
		return 500, nil, err
	}

	item := db.NewItem()
	item.SetTitleAndTags(params.Title)
	item.Content = params.Content
	item.Brhub = b
	item.Type = params.Type
	// validate link
	item.Link = params.Link
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
		return 422, nil, errors.New("invalid item")
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
		return 422, nil, errors.New("invalid item")
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
func AllBrhubs(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	all, err := session.DB().AllBrhubs()
	if err != nil {
		status = 500
	}
	return status, all, err
}

func CreateBrhub(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	var params brhubp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		return 400, nil, errors.New(http.StatusText(400))
	}

	session := appCtx.SessionClone()
	defer session.Close()

	if exists, err := session.DB().BrhubExists(params.Name); err != nil {
		return 500, nil, err
	} else if exists {
		return 422, nil, fmt.Errorf("brhub %s already exists", params.Name)
	}

	b := &db.Brhub{
		Id:       bson.NewObjectId(),
		Name:     params.Name,
		ColorHex: colorful.HappyColor().Hex(),
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
			return 422, nil, errors.New("invalid parent")
		}
		params.parent = bson.ObjectIdHex(params.Parent)

		_, err := session.DB().Comment(params.parent)
		if err != nil {
			if err == mgo.ErrNotFound {
				return 422, nil, errors.New("invalid parent")
			}
			return 500, nil, errors.New(http.StatusText(500))
		}
	}

	_, err := session.DB().Item(params.Item)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 422, nil, errors.New("invalid item")
		}
		return 500, nil, errors.New(http.StatusText(500))
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
	} else {
		err = session.DB().IncrCommentCount(params.Item)
		if err != nil {
			status = 500
		}
	}

	return status, comment, err
}

func Search(appCtx *Context, c web.C, r *http.Request) (int, interface{}, error) {
	session := appCtx.SessionClone()
	defer session.Close()

	term := c.URLParams["term"]
	res := make(map[string]*db.SearchResult)
	funcs := []db.SearchFunc{
		db.SearchThemes,
		db.SearchItems,
	}

	// TODO: bm goroutines
	for _, fn := range funcs {
		sres, err := fn(session.DB(), term)
		if err != nil {
			return 500, nil, err
		}
		if sres.Len > 0 {
			res[sres.Name] = sres
		}
	}
	return 200, map[string]interface{}{
		"results": res,
	}, nil
}
