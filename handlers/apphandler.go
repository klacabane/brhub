package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/klacabane/brhub/db"
	"github.com/zenazn/goji/web"
)

type Context struct {
	Session *db.Session
}

func (ctx *Context) SessionClone() *db.Session {
	return ctx.Session.Clone()
}

type AppHandler struct {
	AppCtx *Context
	H      func(*Context, web.C, http.ResponseWriter, *http.Request) (int, interface{}, error)
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.ServeHTTPC(web.C{}, w, r)
}

func (ah AppHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, data, err := ah.H(ah.AppCtx, c, w, r)
	if err != nil {
		w.WriteHeader(status)
		w.Write(errJson(err.Error()))
		return
	}

	if status == 304 {
		w.WriteHeader(304)
		return
	}

	js, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		w.Write(errJson(http.StatusText(500)))
		return
	}
	w.WriteHeader(status)
	w.Write(js)
}

func errJson(msg string) []byte {
	errstr := fmt.Sprintf("{\"msg\":\"%s\"}", msg)
	return []byte(errstr)
}
