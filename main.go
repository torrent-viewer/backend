package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/torrent-viewer/backend/database"
	"github.com/torrent-viewer/backend/resources/show"
	"github.com/torrent-viewer/backend/router"
)

func main() {
	if err := database.Init("mysql", "torrentviewer:torrentviewer@tcp(127.0.0.1:3306)/tv?charset=utf8&parseTime=True&loc=Local"); err != nil {
		log.Fatal(err)
	}
	database.Conn.AutoMigrate(&show.Show{})
	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(handlers.CORS())
	acceptedTypes := []string{
		"application/vnd.api+json",
	}
	r.Use(router.ContentTypeMiddleware(acceptedTypes))
	show.RegisterHandlers(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
