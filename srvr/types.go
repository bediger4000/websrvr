package srvr

import (
	"net/http"
)

// Srvr holds all info that's used across instances
// of this web application.
type Srvr struct {
	Router *http.ServeMux
	Debug  bool
}
