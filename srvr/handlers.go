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
				multiPart = true
				break
			}
		}

		if multiPart {
			if err := r.ParseMultipartForm(10 * 1024 * 1024); err == nil {
				if r.MultipartForm != nil {
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
							fmt.Printf("\tfilename %q\n", value.Filename)
							f, h, e := r.FormFile(field)
							if e != nil {
								fmt.Printf("Problem on r.FormFile(%q): %v\n", field, e)
								continue
							}
							fmt.Printf("\tSize: %d\n", h.Size)
							fmt.Printf("\tMIME type(s):\n")
							for mimekey, mimevalue := range h.Header {
								fmt.Printf("\t\t%q:%q\n", mimekey, mimevalue)
							}
							buf := make([]byte, 1024)
							sum := 0
							n := 0
							var rerr error
							for n, rerr = f.Read(buf); n > 0 && err == nil; n, err = f.Read(buf) {
								fmt.Printf("%s", string(buf))
								sum += n
							}
							if rerr != nil {
								fmt.Printf("Final error: %v\n", err)
							}
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
