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
		}

		info := LogEntry{
			ReceptionTime: time.Now(),
			Method:        r.Method,
			URL:           r.URL.String(),
			UserAgent:     r.UserAgent(),
			RequestURI:    r.RequestURI,
			Protocol:      r.Proto,
			ContentLength: r.ContentLength,
			Host:          r.Host,
			Remote:        r.RemoteAddr,
		}

		if len(r.TransferEncoding) > 0 {
			info.Encoding = r.TransferEncoding
		}

		hdr := r.Header

		for key, values := range hdr {
			for _, value := range values {
				nvp := &NameValuePair{
					Name:  key,
					Value: value,
				}
				info.Headers = append(info.Headers, nvp)
			}
		}

		s.Data(&info)

		// Return request
		w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}
