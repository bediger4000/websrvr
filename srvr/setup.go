package srvr

import (
	"fmt"
	"os"
	"sync"
)

// Setup prepares Srvr struct for use once all the command line
// flags data gets set in Srvr elements.
func (s *Srvr) Setup() {

	s.logMu = &sync.Mutex{}

	s.LogDescriptor = os.Stderr
	if s.Logfile != "" && s.Logfile != "stderr" && s.Logfile != "-" {
		if logFileName, err := findLogFile(s.Logfile); err == nil {
			if s.Logfile != logFileName {
				if err := os.Rename(s.Logfile, logFileName); err != nil {
					fmt.Fprintf(os.Stderr, "problem moving %q to %q: %v\n", s.Logfile, logFileName, err)
				}
			}
			if fd, err := os.OpenFile(s.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err != nil {
				s.Infof("problem opening log file %q: %v", s.Logfile, err)
				s.Infof("logging to stderr")
			} else {
				s.LogDescriptor = fd
				s.Debugf("logging to file %q", s.Logfile)
			}
		} else {
			fmt.Fprintf(os.Stderr, "problem finding logfile %q: %v\n", s.Logfile, err)
		}
	} else {
		s.Logfile = "stderr"
		s.Infof("logging stderr")
	}

	s.data = StartData(s)

	if s.Address == "" {
		s.Address = fmt.Sprintf(":%s", s.Port)
	}
	s.Infof("Listening on TCP address %q", s.Address)

	if err := os.MkdirAll(s.DownloadDir, 0o777); err != nil {
		s.Infof("problem with download directory %s/", s.DownloadDir, err)
		s.Infof("not saving downloaded files")
		s.DownloadDir = ""
	} else {
		s.Infof("Any downloaded files end up in %s/", s.DownloadDir)
	}
}

func findLogFile(logFileName string) (string, error) {
	// check logFileName
	if _, err := os.Stat(logFileName); err != nil {
		if os.IsNotExist(err) {
			// logFileName exists, check logFileName.0
			nlfn := fmt.Sprintf("%s.0", logFileName)
			if _, err := os.Stat(nlfn); err != nil {
				if os.IsNotExist(err) {
					return logFileName, nil
				}
			}
		} else if !os.IsExist(err) {
			// there's a larger problem
			return "", err
		}
	}

	for i := 0; i < 300; i++ {
		nlfn := fmt.Sprintf("%s.%d", logFileName, i)
		if _, err := os.Stat(nlfn); err != nil {
			if os.IsNotExist(err) {
				// logFileName.N exists, check logFileName.(N+1)
				nlfn2 := fmt.Sprintf("%s.%d", logFileName, i+1)
				if _, err := os.Stat(nlfn2); err != nil {
					if os.IsNotExist(err) {
						return nlfn, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("probably more than 300 versions of %s", logFileName)
}
