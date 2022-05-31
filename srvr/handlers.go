package srvr

import (
	"bytes"
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

		logentry := saveData(s, r)

		gmt, err := time.LoadLocation("GMT")
		if err != nil {
			s.Infof("%v", err)
		}

		// w.Header()["Server"] = []string{"Apache/2.4.53 (Unix) PHP/8.1.4"}
		w.Header()["Server"] = []string{"Server: Apache/2.4.51 (Unix) PHP/5.6.40"}
		w.Header()["Date"] = []string{time.Now().In(gmt).Format(time.RFC1123)}
		/*
			Last-Modified: Thu, 28 May 2020 12:25:56 GMT

			ETag: "82e-5a6b46e2a93c2"
			https://httpd.apache.org/docs/2.2/mod/core.html#FileETag

			Accept-Ranges: bytes
		*/

		if r.URL.String() == `/robots.txt` {
			handleRobotsTxt(w, r)
			return
		}
		if strings.HasSuffix(r.URL.String(), "favicon.ico") {
			handleFaviconIco(w, r)
			return
		}

		if r.Method == "HEAD" {
			handleHead(w, r)
			return
		}

		if strings.HasPrefix(r.URL.String(), "HNAP") {
			handleHnap(w, r)
			s.Infof("HNAP handled")
			return
		}

		// How about a zip bomb for Accept-Encoding: gzip, deflate?
		// What other well-known files?
		// sitemap.xml
		// xmlrpc.ph
		// wp-login.php
		if wsoRequest(r, logentry) {
			handleWso(w, r, logentry)
			s.Infof("WSO handled")
			return
		}

		// Return request
		fmt.Fprintf(w, indexHTML, time.Now().Format(time.RFC3339))
	}
}

func saveData(s *Srvr, r *http.Request) *LogEntry {
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
		// Save Body bytes to a file
		buffer, err := io.ReadAll(r.Body)
		if err != nil {
			s.Infof("problem reading %d body bytes: %v", r.ContentLength, err)
		} else {
			hash := sha256.Sum256(buffer)
			// Identical buffers from different POST requests end up with same name.
			// That's OK: keep from using local disk space to store dupes.
			localFileName := fmt.Sprintf("%s/%s", s.DownloadDir, hex.EncodeToString(hash[:]))
			fout, err := os.Create(localFileName)
			defer fout.Close()
			if err != nil {
				s.Infof("problem creating %s: %v", localFileName, err)
			} else {
				n, err := fout.Write(buffer)
				if err != nil {
					s.Infof("problem writing %s with %d body bytes: %v", localFileName, len(buffer), err)
				} else if n != len(buffer) {
					s.Infof("problem writing %s wrote %d body bytes, wanted to write %d", localFileName, n, len(buffer))
				}
				info.MultiPartFile = localFileName
			}
		}

		// Since r.Body got read, replace with the buffer created to hold it,
		// so that request.ParseForm or request.ParseMultipartForm
		r.Body = io.NopCloser(bytes.NewBuffer(buffer))
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
		// Save the multipart form data as such: it's worthwhile having the
		// form fields and values seperately from the raw data saved above
		if err := r.ParseMultipartForm(10 * 1024 * 1024); err == nil {
			if r.MultipartForm != nil {

				// name/value pairs
				for key, values := range r.MultipartForm.Value {
					for _, value := range values {
						nvp := &NameValuePair{
							Name:  key,
							Value: value,
						}
						info.Form = append(info.Form, nvp)
					}
				}

				// uploaded files, this is gross.
				for field, fileheaders := range r.MultipartForm.File {
					for _, value := range fileheaders {
						fin, h, e := r.FormFile(field)
						if e != nil {
							s.Infof("Problem on r.FormFile(%q): %v\n", field, e)
							continue
						}
						// hopefully a unique local file name, even though this will probably
						// end up with byte-duplicate local file contents. It's difficult
						// to get the uploaded files' bytes for the hash calc (as above for the
						// raw multi-part data) because http.Request makes the files' data
						// available as a *os.File.
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
				s.Infof("nil multi-part form")
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

	return &info
}
