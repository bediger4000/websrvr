package srvr

import (
	"net/http"
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
	Encoding      []string         `json:"transfer_encoding,omitempty"`
	Headers       []*NameValuePair `json:"headers"`
	Cookies       []*CookieEntry   `json:"cookies"`
	Form          []*NameValuePair `json:"form"`
}

// CookieEntry holds HTTP cookie data
type CookieEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`

	Path       string    `json:"path,omitempty"`        // optional
	Domain     string    `json:"domain,omitempty"`      // optional
	Expires    time.Time `json:"expires,omitempty"`     // optional
	RawExpires string    `json:"raw_expires,omitempty"` // for reading cookies only
	MaxAge     int       `json:"max_age"`
	Secure     bool      `json:"secure"`
	HttpOnly   bool      `json:"http_only"`
	SameSite   int       `json:"same_site"`
	Raw        string    `json:"raw"`
	Unparsed   []string  `json:"unparsed,omitempty"`
}
