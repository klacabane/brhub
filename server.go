package main

import (
	"log"
	"net/http"

	"github.com/klacabane/brhub/db"
	"github.com/klacabane/brhub/handlers"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {
	ctx := &handlers.Context{Db: &db.DB{}}

	mux := web.New()
	mux.Get("/api/timeline", handlers.AppHandler{ctx, handlers.Timeline})

	mux.Handle("/*", http.FileServer(http.Dir("./webapp")))

	log.Println("Listening on port 8000")
	graceful.ListenAndServe(":8000", mux)
}
