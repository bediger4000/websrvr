package srvr

import "net/http"

func handleHead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Last-Modified", "Wed, 11 May 2017 13:51:18 GMT")
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", "36501")
}
