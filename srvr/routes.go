package srvr

// Routes sets the URIs and functions that work them.
// "pattern" string documented here: https://golang.org/pkg/net/http/#ServeMux
func (s *Srvr) Routes() {
	s.Router.HandleFunc("/", s.handleSlash())
}
