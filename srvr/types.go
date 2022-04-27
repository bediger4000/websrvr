package srvr

import (
	"net/http"
	"net/textproto"
	"os"
	"sync"
	"time"
)

// Srvr holds all info that's used across instances
// of this web application.
type Srvr struct {
	Port          string
	Address       string
	Router        *http.ServeMux
	Debug         bool
	Logfile       string
	LogDescriptor *os.File
	logMu         *sync.Mutex
	Datafile      string
	data          chan *LogEntry
	OutputLines   int
	DownloadDir   string
}

// NameValuePair carries HTTP header names and values
type NameValuePair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// LogEntry holds HTTP request data in preparation for
// JSON marshalling
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
	Encoding      []string         `json:"transfer_encoding"`
	Headers       []*NameValuePair `json:"headers"`
	Cookies       []*CookieEntry   `json:"cookies"`
	Form          []*NameValuePair `json:"form"`
	Files         []*FileData      `json:"files"`
}

// FileData holds info on 1 file for data logging,
// plus reference (by name) to saved file
type FileData struct {
	FormField     string               `json:"field"`     // name of upload input field
	Size          int64                `json:"size"`      // size of uploaded file
	FileName      string               `json:"filename"`  // uploaded file name - as sent by remote
	LocalFileName string               `json:"localfile"` // uploaded file name here
	MimeGarbage   textproto.MIMEHeader `json:"mimeinfo"`
}

// CookieEntry holds HTTP cookie data
type CookieEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path,omitempty"`
	Domain     string    `json:"domain,omitempty"`
	Expires    time.Time `json:"expires,omitempty"`
	RawExpires string    `json:"raw_expires,omitempty"`
	MaxAge     int       `json:"max_age,omitempty"`
	Secure     bool      `json:"secure,omitempty"`
	HttpOnly   bool      `json:"http_only,omitempty"`
	SameSite   int       `json:"same_site,omitempty"`
	Raw        string    `json:"raw,omitempty"`
	Unparsed   []string  `json:"unparsed,omitempty"`
}
