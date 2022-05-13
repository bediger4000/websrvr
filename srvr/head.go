package srvr

import (
	"net/http"
	"time"
)

func handleHnap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Last-Modified", "Wed, 11 May 2017 13:51:18 GMT")
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Are-You-Kidding", "do/not/be/stupid")
	time.Sleep(25 * time.Second)
}

func handleHead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Last-Modified", "Wed, 11 May 2017 13:51:18 GMT")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", "36501")
}
