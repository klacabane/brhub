package main

import (
	"log"
	"net/http"

	"labix.org/v2/mgo"

	"github.com/klacabane/brhub/handlers"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	if err = session.DB("brhub").C("posts").EnsureIndexKey("-date"); err != nil {
		log.Fatal(err)
	}

	ctx := &handlers.Context{Session: session}

	mux := web.New()
	mux.Get("/api/timeline/:skip/:limit", handlers.AppHandler{ctx, handlers.Timeline})
	mux.Get("/api/b/:id", handlers.AppHandler{ctx, handlers.Items})

	mux.Post("/api/item", handlers.AppHandler{ctx, handlers.CreateItem})
	mux.Get("/api/item/:id", handlers.AppHandler{ctx, handlers.Item})

	mux.Handle("/*", http.FileServer(http.Dir("./webapp")))

	log.Println("Listening on port 8000")
	graceful.ListenAndServe(":8000", mux)
}
