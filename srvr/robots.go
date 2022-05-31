package srvr

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
)

//go:embed "robots.txt"
var robotsTxt string

func handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Last-Modified", "Wed, 11 May 2017 13:51:18 GMT")
	w.Header().Set("ETag", "82e-5a6b46e2a93c2")

	if r.Method == "HEAD" {
		// is this all there is for a HEAD?
		w.Header().Set("Content-Length", strconv.Itoa(len(robotsTxt)))
		return
	}
	fmt.Fprintf(w, robotsTxt)
}
