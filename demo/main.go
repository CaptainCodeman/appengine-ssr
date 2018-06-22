package main

import (
	"time"

	"net/http"

	"google.golang.org/appengine"

	"github.com/captaincodeman/appengine-ssr"
)

func main() {
	mw := ssr.NewSSR("https://pptraas.com",
		ssr.Memcache("ssr:"),
		ssr.Expiration(time.Hour*24*365), // hey, it's a demo
	)

	handler := http.HandlerFunc(handle)
	http.Handle("/", mw.Middleware(handler))
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
