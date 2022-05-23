package srvr

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func wsoRequest(r *http.Request, le *LogEntry) bool {
	if strings.HasSuffix(r.URL.String(), "wso.php") {
		return true
	}
	if strings.HasSuffix(r.URL.String(), "twentytwelve/404.php") {
		return true
	}

	// a, ajax, c, charset, pass, p1, p2, p3 parameters
	form := r.Form
	if len(form["a"]) > 0 {
		wsoRequest := false
		for _, value := range form["a"] {
			switch strings.ToLower(value) {
			case "filesman", "network", "rc", "secinfo", "filestools", "console", "php", "stringtools", "safemode", "selfremove", "bruteforce":
				wsoRequest = true
				break
			}
		}
		return wsoRequest
	}
	// What about just a "pass" form value, nothing else?
	if len(form["pass"]) > 0 {
		// presence of the Vigilant Cookie would be a dead giveaway
	}

	// what about the 4.x WSO series requests?

	return false
}

func handleWso(w http.ResponseWriter, r *http.Request, le *LogEntry) {
	fmt.Fprintf(w, "<pre align=center><form method=post>Password: <input type=password name=pass><input type=submit value='>>'></form></pre>")
	time.Sleep(20 * time.Second)
}
