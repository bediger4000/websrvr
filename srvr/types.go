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
	Form          []*NameValuePair `json:"form"`
}
