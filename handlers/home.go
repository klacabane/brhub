package handlers

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func Timeline(appCtx *Context, c web.C, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	return 200, struct{ Items []string }{[]string{"foo", "foo"}}, nil
}
