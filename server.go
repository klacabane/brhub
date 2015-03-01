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
	session := db.MainSession()

	ctx := &handlers.Context{Session: session}

	mux := web.New()
	mux.Get("/api/timeline/:skip/:limit", handlers.AppHandler{ctx, handlers.Timeline})

	mux.Post("/api/b", handlers.AppHandler{ctx, handlers.CreateBrhub})
	mux.Get("/api/b/:id/:skip/:limit", handlers.AppHandler{ctx, handlers.Items})

	mux.Post("/api/items", handlers.AppHandler{ctx, handlers.CreateItem})
	mux.Get("/api/items/:id", handlers.AppHandler{ctx, handlers.Item})

	mux.Post("/api/comments", handlers.AppHandler{ctx, handlers.CreateComment})

	mux.Handle("/*", http.FileServer(http.Dir("./webapp")))

	log.Println("Listening on port 8000")
	graceful.ListenAndServe(":8000", mux)
}
