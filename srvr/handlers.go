package srvr

import (
	"fmt"
	"net/http"
	"time"
)

var indexHTML = `<html>
<head>
</head>
<body>
<p>Time is %s</p>
</body>
</html>
`

func (s *Srvr) handleSlash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		if err := r.ParseForm(); err == nil {
			if len(r.Form) > 0 {
				for key, values := range r.Form {
					for _, value := range values {
						nvp := &NameValuePair{
							Name:  key,
							Value: value,
						}
						info.Form = append(info.Form, nvp)
					}
				}
			}
		} else {
			s.Infof("http.Request.ParseForm(): %v", err)
		}

		contentTypes := r.Form["Content-Type"]
		for _, ct := range contentTypes {
			if ct == "multipart/form-data" {
				if err := r.ParseMultipartForm(10 * 1024 * 1024); err == nil {
					fmt.Printf("multi-part form %T, %v\n", r.MultipartForm, r.MultipartForm)
					if r.MultipartForm != nil {
						fmt.Printf("multi-part form\n")
						for key, value := range r.MultipartForm.Value {
							fmt.Printf("%q: %q\n", key, value)
						}
						/*
							type Form struct {
								Value map[string][]string
								File  map[string][]*FileHeader
							}
						*/
					} else {
					}
				} else {
					s.Infof("http.Request.MultipartParseForm(): %v", err)
				}
			}
		}

		s.data <- &info

		// Return request
		w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}
