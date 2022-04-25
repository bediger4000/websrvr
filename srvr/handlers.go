package srvr

import (
	"fmt"
	"net/http"
	"strings"
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

		for key, values := range r.Header {
			for _, value := range values {
				nvp := &NameValuePair{
					Name:  key,
					Value: value,
				}
				info.Headers = append(info.Headers, nvp)
			}
		}

		for _, c := range r.Cookies() {
			ce := &CookieEntry{
				Name:       c.Name,
				Value:      c.Value,
				Path:       c.Path,
				Expires:    c.Expires,
				RawExpires: c.RawExpires,
				MaxAge:     c.MaxAge,
				Secure:     c.Secure,
				HttpOnly:   c.HttpOnly,
				SameSite:   int(c.SameSite),
				Raw:        c.Raw,
				Unparsed:   c.Unparsed,
			}
			info.Cookies = append(info.Cookies, ce)
		}

		multiPart := false
		contentTypes := r.Header["Content-Type"]
		for _, ct := range contentTypes {
			if strings.HasPrefix(ct, "multipart/form-data") {
				fmt.Printf("multipart/form-data 0\n")
				multiPart = true
				break
			}
		}

		if multiPart {
			fmt.Printf("multipart/form-data 1\n")
			if err := r.ParseMultipartForm(10 * 1024 * 1024); err == nil {
				fmt.Printf("multipart/form-data 2\n")
				if r.MultipartForm != nil {
					fmt.Printf("multipart/form-data 3\n")
					for key, values := range r.MultipartForm.Value {
						for _, value := range values {
							nvp := &NameValuePair{
								Name:  key,
								Value: value,
							}
							info.Form = append(info.Form, nvp)
						}
					}
					for field, fileheaders := range r.MultipartForm.File {
						for _, value := range fileheaders {
							fmt.Printf("field %q\n", field)
							fmt.Printf("\tfilename %q\n", value.Filename)
						}
					}
					/*
						type FileHeader struct {
							Filename string
							Header   textproto.MIMEHeader
							Size     int64

							content []byte
							tmpfile string
						}

						type Form struct {
							Value map[string][]string
							File  map[string][]*FileHeader
						}
					*/
				} else {
					fmt.Printf("nil multi-part form\n")
				}
			} else {
				s.Infof("http.Request.MultipartParseForm(): %v", err)
			}
		} else {
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
		}

		s.data <- &info

		// Return request
		w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}
