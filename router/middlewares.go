package router

import (
	"log"
	"net/http"
	"time"
)

type logger struct {
	h http.Handler
}

func (l logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.h.ServeHTTP(w, r)
	log.Printf("%10s\t%-50s\t%s\n", r.Method, r.URL.Path, time.Since(start))
}

// LoggingMiddleware logs the requested URL and the time spent in each route
func LoggingMiddleware(handler http.Handler) http.Handler {
	return logger{
		h: handler,
	}
}

type contentType struct {
	h        http.Handler
	accepted []string
}

func (c contentType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, t := range c.accepted {
		if t == r.Header.Get("Content-Type") {
			c.h.ServeHTTP(w, r)
		}
	}
	w.WriteHeader(http.StatusUnsupportedMediaType)
}

// ContentTypeMiddleware restricts the Content-Type that can be requested
func ContentTypeMiddleware(accepted []string) Middleware {
	return func(handler http.Handler) http.Handler {
		return contentType{
			h:        handler,
			accepted: accepted,
		}
	}
}
