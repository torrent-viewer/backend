package router

import (
	"log"
	"net/http"
	"time"
	"regexp"
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
			return
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

type Guard func(r *http.Request) bool

type firewall struct {
	only   []*regexp.Regexp
	except []*regexp.Regexp
	guard  Guard
	h      http.Handler
}

type FirewallConfig struct {
	Only   []string
	Except []string
	Guard  Guard
}

func (fw firewall) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authenticated := false
	if len(fw.only) > 0 {
		for _, pattern := range fw.only {
			if pattern.MatchString(r.URL.Path) == true {
				authenticated = fw.guard(r)
				if authenticated {
					break;
				}
			}
		}
		if authenticated == false {
			w.WriteHeader(401)
			return
		}
	} else if len(fw.except) > 0 {
		for _, pattern := range fw.except {
			if pattern.MatchString(r.URL.Path) == false {
				authenticated = fw.guard(r)
				if authenticated {
					break;
				}
			}
		}
		if authenticated == false {
			w.WriteHeader(401)
			return
		}
	}
	fw.h.ServeHTTP(w, r)
}

func firewallCompileSlice(patterns []string) []*regexp.Regexp {
	compiled := make([]*regexp.Regexp, len(patterns), len(patterns))
	for i, pattern := range patterns {
		compiled[i] = regexp.MustCompile(pattern)
	}
	return compiled
}

func FirewallMiddleware(config FirewallConfig) Middleware {
	onlyCompiled := firewallCompileSlice(config.Only)
	exceptCompiled := firewallCompileSlice(config.Except)
	return func(handler http.Handler) http.Handler {
		return firewall{
			only:   onlyCompiled,
			except: exceptCompiled,
			h:      handler,
			guard:  config.Guard,
		}
	}
}
