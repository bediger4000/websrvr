package srvr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
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

		saveData(s, r)

		w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}

		if r.URL.String() == `/robots.txt` {
			handleRobotsTxt(w, r)
			return
		}
		if strings.HasSuffix(r.URL.String(), "favicon.ico") {
			handleFaviconIco(w, r)
			return
		}
		if wsoRequest(r) {
			handleWso(w, r)
			return
		}

		// Return request
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}

func saveData(s *Srvr, r *http.Request) {
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
	} else {
		info.Encoding = make([]string, 0)
	}

	info.Headers = make([]*NameValuePair, 0)
	for key, values := range r.Header {
		for _, value := range values {
			nvp := &NameValuePair{
				Name:  key,
				Value: value,
			}
			info.Headers = append(info.Headers, nvp)
		}
	}

	// Set some struct elements to a zero-length slice to prevent
	// JSON output from being "slicename": null
	info.Cookies = make([]*CookieEntry, 0)
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

	info.Files = make([]*FileData, 0)
	info.Form = make([]*NameValuePair, 0)

	if r.ContentLength > 0 {
		s.Infof("content length %d > 0", r.ContentLength)
		buffer := make([]byte, r.ContentLength)
		n, err := r.Body.Read(buffer)
		if f, ok := r.Body.(*os.File); ok {
			f.Seek(0, os.SEEK_SET)
		}
		if err != nil {
			s.Infof("reading %d body bytes: %v", r.ContentLength, err)
		} else {
			s.Infof("read %d body bytes", n)
			if n != int(r.ContentLength) {
				s.Infof("read %d body bytes, wanted to read %d", n, r.ContentLength)
			}
			hash := sha256.Sum256(buffer)
			localFileName := fmt.Sprintf("%s/%s", s.DownloadDir, hex.EncodeToString(hash[:]))
			s.Infof("local file name ", localFileName)
			fout, err := os.Create(localFileName)
			if err != nil {
				s.Infof("creating %s: %v", localFileName, err)
			} else {
				n, err := fout.Write(buffer)
				if err != nil {
					s.Infof("filling %s with $d bytes: %v", localFileName, len(buffer), err)
				} else if n != len(buffer) {
					s.Infof("filling %s wrote %d bytes, wanted to write %d", localFileName, n, len(buffer))
				}
				fd := &FileData{
					FormField:     "whole body",
					Size:          int64(n),
					LocalFileName: localFileName,
				}
				info.Files = append(info.Files, fd)
			}
			fout.Close()
		}
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
						// fmt.Printf("\tfilename %q\n", value.Filename)
						fin, h, e := r.FormFile(field)
						if e != nil {
							s.Infof("Problem on r.FormFile(%q): %v\n", field, e)
							continue
						}
						hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", field, value.Filename, time.Now().Format(time.RFC3339))))
						localFileName := fmt.Sprintf("%s/%s", s.DownloadDir, hex.EncodeToString(hash[:]))
						fout, err := os.Create(localFileName)
						skipOutput := false
						if err != nil {
							s.Infof("creating %s: %v", localFileName, err)
							localFileName = "error creating"
							skipOutput = true
						}
						fd := &FileData{
							FormField:     field,
							Size:          h.Size,
							FileName:      h.Filename,
							LocalFileName: localFileName,
							MimeGarbage:   h.Header,
						}
						if !skipOutput {
							n, err := io.Copy(fout, fin)
							if err != nil {
								s.Infof("filling %s: %v", localFileName, err)
							}
							if n != h.Size {
								s.Infof("filling %s wrote %d bytes, wanted to write %d", localFileName, n, h.Size)
							}
						}
						fin.Close()
						fout.Close()
						info.Files = append(info.Files, fd)
					}
				}
			} else {
				s.Infof("nil multi-part form\n")
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
}

func handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, `User-agent: Googlebot
Disallow: /wp/

User-agent: *
Allow: /

User-agent: *
Disallow: /porn

User-agent: *
Disallow: /flapjackattack
`,
	)
}

var favicon = []byte{
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10, 0x10, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x28, 0x01, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x80,
	0x00, 0x00, 0x80, 0x80, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x80, 0x00,
	0x80, 0x00, 0x00, 0x80, 0x80, 0x00, 0xc0, 0xc0, 0xc0, 0x00, 0x80, 0x80,
	0x80, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0xff,
	0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0x00, 0xff,
	0xff, 0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x0f, 0xf0,
	0x00, 0x00, 0x00, 0x00, 0x0f, 0xf0, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00,
	0xf0, 0x00, 0x00, 0x00, 0xf0, 0xff, 0xff, 0x0f, 0x00, 0x00, 0x00, 0x00,
	0x0f, 0xf0, 0x0f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x0f, 0xff, 0xff, 0xf0, 0x00, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x0f, 0xf0, 0xff,
	0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x0f, 0x00, 0xff, 0xff, 0x00,
	0xf0, 0x00, 0x0f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x0f, 0xf0, 0x00, 0xf0,
	0x00, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a,
}

func handleFaviconIco(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(favicon)
}
