package srvr

import (
	"fmt"
	"net/http"
	"strings"
)

func wsoRequest(r *http.Request) bool {
	if strings.HasSuffix(r.URL.String(), "wso.php") {
		return true
	}
	if strings.HasSuffix(r.URL.String(), "twentytwelve/404.php") {
		return true
	}

	// a, c, p1, p2 parameters in request

	// what about the 4.x WSO series requests?

	return false
}

func handleWso(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<pre align=center><form method=post>Password: <input type=password name=pass><input type=submit value='>>'></form></pre>")
}
