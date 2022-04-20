package srvr

import (
	"fmt"
	"os"
)

// Setup prepares Srvr struct for use once all the command line
// flags data gets set in Srvr elements.
func (s *Srvr) Setup() {

	s.LogDescriptor = os.Stderr
	if s.Logfile != "" && s.Logfile != "stderr" && s.Logfile != "-" {
		if fd, err := os.OpenFile(s.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err != nil {
			s.Infof("problem opening log file %q: %v", s.Logfile, err)
			s.Infof("logging to stderr")
		} else {
			s.LogDescriptor = fd
			s.Debugf("logging to file %q", s.Logfile)
		}
	} else {
		s.Logfile = "stderr"
		s.Infof("logging stderr")
	}

	s.DataDescriptor = os.Stdout
	if s.Datafile != "" {
		if fd, err := os.OpenFile(s.Datafile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err != nil {
			s.Infof("problem opening data file %q: %v", s.Logfile, err)
			s.Infof("data written to stdout")
		} else {
			s.DataDescriptor = fd
			s.Infof("data sent to file %q", s.Datafile)
		}
	} else {
		s.Datafile = "stdout"
		s.Infof("data written to stdout")
	}

	if s.Address == "" {
		s.Address = fmt.Sprintf(":%s", s.Port)
	}
	s.Infof("Listening on TCP address %q", s.Address)
}
