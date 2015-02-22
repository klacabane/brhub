package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"labix.org/v2/mgo"

	"github.com/klacabane/brhub/db"
	"github.com/zenazn/goji/web"
)

type Context struct {
	Session *mgo.Session
}

func (ctx *Context) SessionClone() *db.Session {
	return &db.Session{ctx.Session.Clone()}
}

type AppHandler struct {
	AppCtx *Context
	H      func(*Context, web.C, http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.ServeHTTPC(web.C{}, w, r)
}

func (ah AppHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	status, data, err := ah.H(ah.AppCtx, c, w, r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err != nil {
		w.Write(errJson(err.Error()))
		return
	}

	js, err := json.Marshal(data)
	if err != nil {
		w.Write(errJson(http.StatusText(500)))
		return
	}
	w.Write(js)
}

func errJson(msg string) []byte {
	errstr := fmt.Sprintf("{\"msg\":\"%s\"}", msg)
	return []byte(errstr)
}
