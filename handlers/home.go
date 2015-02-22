package handlers

import (
	"fmt"
	"net/http"

	"labix.org/v2/mgo/bson"

	"github.com/zenazn/goji/web"
)

func Timeline(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	posts, err := session.DB().Timeline(bson.NewObjectId() /*userid*/)
	if err != nil {
		status = 500
	}
	return status, posts, err
}

func Posts(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	if !bson.IsObjectIdHex(c.URLParams["id"]) {
		return 400, nil, fmt.Errorf("invalid id")
	}
	id := bson.ObjectIdHex(c.URLParams["id"])

	session := appCtx.SessionClone()
	defer session.Close()

	status := 200
	posts, err := session.DB().Posts(id)
	if err != nil {
		status = 500
	}
	return status, posts, err
}
