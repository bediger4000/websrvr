package srvr

import (
	"fmt"
	"net/http"
	"time"
)

var indexHTML = `
<html>
<head>
</head>
<body>
<p>Time is %s</p>
</body>
</html>
`

func (s *Srvr) handleSlash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleSlash closure\n")
			defer fmt.Printf("Exit handleSlash closure\n")

			fmt.Printf("Method: %q\n", r.Method)
			fmt.Printf("URL: %+v\n", r.URL)
			fmt.Printf("Proto: %q\n", r.Proto)
			fmt.Printf("ContentLength: %d\n", r.ContentLength)
			fmt.Printf("Host: %q\n", r.Host)
			fmt.Printf("Remote Address: %q\n", r.RemoteAddr)
			fmt.Printf("RequestURI: %q\n", r.RemoteAddr)
			if len(r.TransferEncoding) > 0 {
				fmt.Printf("TransferEncoding: %v\n", r.TransferEncoding)
			}

			hdr := r.Header
			fmt.Printf("Found %d request headers\n", len(hdr))

			for key, values := range hdr {
				fmt.Printf("Header: %q\n", key)
				for _, value := range values {
					fmt.Printf("\t%q\n", value)
				}
			}
		}
		hdr := w.Header()
		hdr["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}
