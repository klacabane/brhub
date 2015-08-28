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
	session := db.MainSession("localhost")
	ctx := &handlers.Context{Session: session}

	api := web.New()
	api.Use(handlers.ValidateToken)
	api.Get("/api/timeline/:skip/:limit", handlers.AppHandler{ctx, handlers.Timeline})

	api.Post("/api/b/", handlers.AppHandler{ctx, handlers.CreateTheme})
	api.Get("/api/b/", handlers.AppHandler{ctx, handlers.AllThemes})
	api.Get("/api/b/:name/:skip/:limit", handlers.AppHandler{ctx, handlers.Items})

	api.Post("/api/items/", handlers.AppHandler{ctx, handlers.CreateItem})
	api.Get("/api/items/:id", handlers.AppHandler{ctx, handlers.Item})
	api.Patch("/api/items/:id/upvote", handlers.AppHandler{ctx, handlers.Upvote})

	api.Post("/api/comments/", handlers.AppHandler{ctx, handlers.CreateComment})

	api.Patch("/api/users/stars", handlers.AppHandler{ctx, handlers.Star})

	api.Get("/api/search/:term", handlers.AppHandler{ctx, handlers.Search})

	mux := web.New()
	mux.Post("/auth", handlers.AppHandler{ctx, handlers.Auth})
	mux.Handle("/api/*", api)
	mux.Handle("/*", http.FileServer(http.Dir("./webapp")))

	log.Println("Listening on port 8000")
	graceful.ListenAndServe(":8000", mux)
}
