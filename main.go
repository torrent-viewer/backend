package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/torrent-viewer/backend/datastore"
	"github.com/torrent-viewer/backend/resources/show"
	"github.com/torrent-viewer/backend/router"
)

func BasicAuth(r *http.Request) bool {
	return true
}

func main() {
	dbDriver := os.Getenv("TV_DB_DRIVER")
	dbUser := os.Getenv("TV_DB_USER")
	dbPassword := os.Getenv("TV_DB_PASSWORD")
	dbHost := os.Getenv("TV_DB_HOST")
	dbPort := os.Getenv("TV_DB_PORT")
	dbBase := os.Getenv("TV_DB_BASE")
	log.Printf("Connection to %s database: %s:%s@%s:%s/%s", dbDriver, dbUser, dbPassword, dbHost, dbPort, dbBase)
	err := datastore.Init(dbDriver, dbUser, dbPassword, dbHost, dbPort, dbBase)
	if err != nil {
		for {
			log.Println("Could not connect to database\n", err, "\nRetrying in 1 second...")
			time.Sleep(1000 * time.Millisecond)
			err = datastore.Init(dbDriver, dbUser, dbPassword, dbHost, dbPort, dbBase)
			if err == nil {
				break
			}
		}
	}
	datastore.Conn.AutoMigrate(&show.Show{})
	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(handlers.CORS())
	acceptedTypes := []string{
		"application/vnd.api+json",
		"application/vnd.api+json; charset=UTF-8",
		"application/vnd.api+json; charset=utf-8",
	}
	r.Use(router.ContentTypeMiddleware(acceptedTypes))
	r.Use(router.FirewallMiddleware(router.FirewallConfig{
		Guard: BasicAuth,
		Only: []string{"^/shows"},
	}))
	r.AddResource("shows", show.ShowResource{})
	log.Fatal(http.ListenAndServe(":8080", r))
}
