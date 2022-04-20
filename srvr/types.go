package srvr

import (
	"net/http"
	"os"
)

// Srvr holds all info that's used across instances
// of this web application.
type Srvr struct {
	Router         *http.ServeMux
	Debug          bool
	Logfile        string
	LogDescriptor  *os.File
	Datafile       string
	DataDescriptor *os.File
}
