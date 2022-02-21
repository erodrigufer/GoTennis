package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// store all flag-parseable config values in this struct
type configValues struct {
	addr string
	//StaticDir string
}

// handle application-wide dependencies in this struct
// this dependencies are then 'injected' to the different handlers,
// by defining the handlers as methods to this struct
// handling the dependencies in this ways makes the code more easy to unit-test,
// just defining these dependencies as global would not make the code easier to
// unit-test
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Define default HOST and PORT, in case flag is not present
	DEFAULT_SERVICE := "localhost:4000"

	cfg := new(configValues)
	flag.StringVar(&cfg.addr, "addr", DEFAULT_SERVICE, "Server's listening address")
	flag.Parse()

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize an instance of application containing the application-wide
	// dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Get the mux from the method at routing.go
	mux := app.routes()
	// Initialize a new http.Server struct.
	// Use errorLog for errors instead of default option
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Start a TCP web server listening on addr
	// If ListenAndServe() returns an error we use the log.Fatal()
	// function to log the error message and exit.
	// ListenAndServe handles all accepted clients concurrently in goroutines
	infoLog.Printf("Starting server at %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Eduardo Rodriguez @erodrigufer (c) 2022
