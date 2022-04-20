package srvr

import (
	"os"
)

// Setup prepares Srvr struct for use once all the command line
// flags data gets set in Srvr elements.
func (s *Srvr) Setup() {
	s.LogDescriptor = os.Stderr
	if s.Logfile != "" {
		if fd, err := os.OpenFile(s.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err != nil {
			s.Infof("problem opening log file %q: %v", s.Logfile, err)
			s.Infof("logging to stderr")
		} else {
			s.LogDescriptor = fd
			s.Debugf("logging to file %q", s.Logfile)
		}
	} else {
		s.Logfile = "stderr"
	}
}
