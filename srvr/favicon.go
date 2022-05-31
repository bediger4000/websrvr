package srvr

import (
	_ "embed"
	"net/http"
)

//go:embed "favicon.ico"
var favicon []byte

func handleFaviconIco(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Last-Modified", "Wed, 11 May 2017 13:51:18 GMT")
	w.Header().Set("ETag", "28e-a6b46e2a93c25")
	w.Write(favicon)
}
