package srvr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogEntry struct {
	ReceptionTime time.Time        `json:"recpt_time"`
	Method        string           `json:"method"`
	URL           string           `json:"url"`
	UserAgent     string           `json:"user_agent"`
	RequestURI    string           `json:"request_uri"`
	Protocol      string           `json:"proto"`
	ContentLength int64            `json:"content_len"`
	Host          string           `json:"host_addr"`
	Remote        string           `json:"remote_addr"`
	Headers       []*NameValuePair `json:"headers"`
	Encoding      []string         `json:"transfer_encoding,omitempty"`
}

func (s *Srvr) handleSlash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			fmt.Printf("Enter handleSlash closure\n")
			defer fmt.Printf("Exit handleSlash closure\n")
		}

		info := &LogEntry{
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

		if buf, err := json.Marshal(&info); err != nil {
			log.Printf("marshalling log JSON: %v\n", err)
		} else {
			n, err := os.Stderr.Write(buf)
			if n != len(buf) {
				log.Printf("wrote %d bytes of JSON, should have written %d\n", n, len(buf))
			}
			if err != nil {
				log.Printf("writing %d bytes log JSON: %v\n", len(buf), err)
			}
		}

		// Return request
		w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}
