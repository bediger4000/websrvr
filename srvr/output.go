package srvr

import (
	"fmt"
	"log"
	"time"
)

// Infof always writes format and args to log file.
func (s *Srvr) Infof(format string, a ...any) {
	_, err := fmt.Fprintf(s.LogDescriptor, fmt.Sprintf("%v\t%s\n", time.Now().Format(time.RFC3339Nano), format), a...)
	if err != nil {
		log.Printf("trying to write info to log file: %v\n", err)
	}
}

// Debugf writes format and args to log file if verbose output flag set
func (s *Srvr) Debugf(format string, a ...any) {
	if !s.Debug {
		return
	}
	s.Infof(format, a...)
}
