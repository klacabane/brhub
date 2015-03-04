package handlers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/klacabane/brhub/db"
	"github.com/zenazn/goji/web"
	"labix.org/v2/mgo/bson"
)

func ValidateToken(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if c.Env == nil {
			c.Env = map[string]interface{}{}
		}

		author, err := parseToken(r.Header.Get("X-token"))
		if err != nil {
			w.WriteHeader(401)
			w.Write(errJson(err.Error()))
			return
		}
		c.Env["user"] = author

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func parseToken(val string) (db.Author, error) {
	token, err := jwt.Parse(val, func(t *jwt.Token) (interface{}, error) {
		return []byte("secretkey"), nil
	})

	if err == nil && token.Valid && bson.IsObjectIdHex(token.Claims["id"].(string)) {
		return db.Author{
			Id:   bson.ObjectIdHex(token.Claims["id"].(string)),
			Name: token.Claims["name"].(string),
		}, nil
	}
	return db.Author{}, fmt.Errorf("invalid token")
}
