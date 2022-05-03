package srvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Output struct {
	C           chan *LogEntry
	Datafile    string
	OutputCount int
	S           *Srvr
	currentFD   *os.File
}

func StartData(s *Srvr) chan *LogEntry {
	o := &Output{
		C:           make(chan *LogEntry, 10),
		Datafile:    s.Datafile,
		S:           s,
		OutputCount: s.OutputLines,
	}

	if s.Datafile != "" {
		o.openDataFile()
	} else {
		o.currentFD = os.Stdout
		o.S.Datafile = "stdout"
		o.S.Infof("data written to stdout")
	}

	go o.WriteData()

	return o.C
}

func (o *Output) openDataFile() {
	if dstDataFile, err := findLogFile(o.Datafile); err == nil {
		if o.Datafile != dstDataFile {
			if err := os.Rename(o.Datafile, dstDataFile); err != nil {
				fmt.Fprintf(os.Stderr, "problem moving %q to %q: %v\n", o.Datafile, dstDataFile, err)
			}
		}
	}
	if fd, err := os.OpenFile(o.Datafile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err != nil {
		o.S.Infof("problem opening data file %q: %v", o.Datafile, err)
		o.S.Infof("data written to stdout")
		o.currentFD = os.Stdout
	} else {
		o.currentFD = fd
		o.S.Infof("data sent to file %q", o.Datafile)
	}
}

func (o *Output) WriteData() {

	output := 0

	for entry := range o.C {
		if buf, err := json.Marshal(entry); err != nil {
			o.S.Infof("marshalling log JSON: %v", err)
		} else {
			buf = append(buf, '\n')
			n, err := o.currentFD.Write(buf)
			if n != len(buf) {
				o.S.Infof("wrote %d bytes of JSON, should have written %d", n, len(buf))
				continue
			}
			if err != nil {
				o.S.Infof("trying to write %d bytes data JSON: %v", len(buf), err)
				continue
			}
			output++
			if o.S.Debug && n == len(buf) && err == nil {
				o.S.Infof("wrote %d bytes of data JSON", n)
			}
		}

		if output > o.OutputCount {
			if o.currentFD == os.Stdout {
				continue
			}
			output = 0
			// Roll file
			o.S.Debugf("rolling log file")
			for i := 0; true; i++ {
				newDatafile := fmt.Sprintf("%s.%d", o.Datafile, i)
				if _, err := os.Stat(newDatafile); err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						err := os.Rename(o.Datafile, newDatafile)
						if err != nil {
							o.S.Infof("renaming %s to %s: %v", o.Datafile, newDatafile)
							break
						}
						o.S.Debugf("newest old logfile %s", newDatafile)
						o.currentFD.Close()
						o.openDataFile()
						break
					}
				}
			}
		}
	}
}
