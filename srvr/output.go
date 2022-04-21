package srvr

import (
	"encoding/json"
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

// Data writes buffer to data fileL
func (s *Srvr) Data(entry *LogEntry) {

	if buf, err := json.Marshal(entry); err != nil {
		s.Infof("marshalling log JSON: %v", err)
	} else {
		n, err := s.DataDescriptor.Write(buf)
		if n != len(buf) {
			s.Infof("wrote %d bytes of JSON, should have written %d", n, len(buf))
		}
		if err != nil {
			s.Infof("trying to write %d bytes data JSON: %v", len(buf), err)
		}
		if s.Debug && n == len(buf) && err == nil {
			s.Infof("wrote %d bytes of data JSON", n)
		}
	}
}
