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

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the root app method as the handler for the "/" URL pattern.
	// instead of using functions, use methods to inject the packet-wide
	// dependencies (like the loggers) in every method without using global
	// variables
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.root)
	mux.HandleFunc("/session/create", app.createSession)
	mux.HandleFunc("/session", app.showSession)

	// Create a handler/fileServer for all files in the static directory
	// Type Dir implements the interface required by FileServer and makes the
	// code portable by using the native file system (which could be different
	// for Windows and other Unix systems)
	// ./ui/static/ will be the root (like a root jail) of the fileServer
	// it will serve files relative to this path. Nonetheless, a security
	// concern is that symlink that points outside the 'jail' can also be
	// followed (check documentation of type Dir)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// the fileServer is now the handler for all URL paths starting with
	// '/static/'
	// http.StripPrefix, will create a new http.Handler that first strips the
	// prefix "/static" from the URL request, and passes the new request to
	// the fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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
