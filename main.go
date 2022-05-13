package main

import (
	"flag"
	"log"
	"net/http"
	"time"
	"websrvr/srvr"
)

func main() {
	portString := flag.String("p", "8080", "TCP port on which to listen")
	addressString := flag.String("a", "", "TCP address on which to listen")
	debug := flag.Bool("v", false, "verbose output per request")
	logFile := flag.String("o", "websrvr.log", "log file")
	dataFile := flag.String("d", "websrvr.data", "JSON data file")
	outputLines := flag.Int("N", 50000, "output lines in data file")
	downloadDir := flag.String("D", "downloads", "downloaded file directory")

	flag.Parse()

	srv := &srvr.Srvr{
		Port:        *portString,
		Address:     *addressString,
		Router:      http.NewServeMux(),
		Debug:       *debug,
		Logfile:     *logFile,
		Datafile:    *dataFile,
		OutputLines: *outputLines,
		DownloadDir: *downloadDir,
	}

	srv.Setup()
	srv.Routes()

	s := &http.Server{
		Addr:    srv.Address,
		Handler: srv.Router,
		// long timeouts to allow us to tarpit some requests
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// signal handling
	// s.Close()
	// s.Shutdown()

	log.Fatal(s.ListenAndServe())
}
