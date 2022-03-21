package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestApplication helper returns an instance of our application struct
// containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(ioutil.Discard, "", 0),
		infoLog:  log.New(ioutil.Discard, "", 0),
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// newTestServer helper initalizes and returns a new instance of the custom
// testServer type as a TLSServer, with a given http.Handler as parameter
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	// initialize a GET request from a test client to a particular path on the
	// testServer
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
