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
	debug := flag.Bool("v", false, "verbose output per request")
	logFile := flag.String("o", "websrvr.log", "log file")
	dataFile := flag.String("d", "websrvr.data", "JSON data file")
	flag.Parse()

	srv := &srvr.Srvr{
		Router:   http.NewServeMux(),
		Debug:    *debug,
		Logfile:  *logFile,
		Datafile: *dataFile,
	}

	srv.Setup()
	srv.Routes()

	s := &http.Server{
		Addr:           ":" + *portString,
		Handler:        srv.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
