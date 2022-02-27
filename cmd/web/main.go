package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/erodrigufer/GoTennis/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql" // the driver's init() function must be
	// run so that it can register itself with "database/sql" nothing else is
	// actually used from this packet, so if no underscore would be present the
	// Go compiler would bring up an error
	"github.com/golangcollege/sessions" // session manager
)

// store all flag-parseable config values in this struct
type configValues struct {
	addr   string
	dsn    string
	secret string
	//StaticDir string
}

// handle application-wide dependencies in this struct
// this dependencies are then 'injected' to the different handlers,
// by defining the handlers as methods to this struct
// handling the dependencies in this ways makes the code more easy to unit-test,
// just defining these dependencies as global would not make the code easier to
// unit-test
type application struct {
	errorLog       *log.Logger                   // error log handler
	infoLog        *log.Logger                   // info log handler
	sessionManager *sessions.Session             // session manager
	session        *mysql.SessionModel           // db for application
	templateCache  map[string]*template.Template // Cache map with html templates
}

func main() {
	// Define default HOST and PORT, in case flag is not present
	DEFAULT_SERVICE := ":4000"

	cfg := new(configValues)
	flag.StringVar(&cfg.addr, "addr", DEFAULT_SERVICE, "Server's listening address")
	// dsn is needed to know how to connect to a db
	// it is composed of ${USERNAME}:${PASSWORD}@/${DB_NAME}?${FLAGS}
	// parseTime=true converts SQL TIME and DATE fields to Go time.Time objects
	flag.StringVar(&cfg.dsn, "dsn", "web:Password1@/goTennis?parseTime=true", "DSN (Data Source Name) for MySQL db")
	// Session secret (a random key) used to encrypt and authenticate session
	// cookies. It should be 32 bytes long
	flag.StringVar(&cfg.secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Session's secret key to encrypt and authenticate session cookies")
	flag.Parse()

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a connection to a db connection pool
	db, err := connectDBpool(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// close connection pool before main() exits
	defer db.Close()

	// Initialize the templates' cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a session manager, pass the secret key and configure the
	// manager so that it always expires after 12 hours
	sessionManager := sessions.New([]byte(*cfg.secret))
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize an instance of application containing the application-wide
	// dependencies
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		session:        &mysql.SessionModel{DB: db},
		sessionManager: sessionManager,
		templateCache:  templateCache,
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// wrapper for sql.Open, return a sql.DB connection pool
// dsn stands for data source name, and is needed to pass the authentication
// information to the DB among other things
// Troubleshooting: if the db service is not correctly active,
// then running `ss -at` will not show the mysql db listening in localhost
// at port `mysql`
func connectDBpool(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// check if connection was established by pinging the db
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Eduardo Rodriguez @erodrigufer (c) 2022
