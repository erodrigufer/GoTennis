package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Initialize a new time.Time object and pass it to the humanDate function
	tm := time.Date(2020, 3, 10, 9, 0, 0, 0, time.UTC)
	expected := "10 Mar 2020 at 09:00"
	hd := humanDate(tm)

	// Check that output from humanDate is in the format that we expect, if it
	// is not, then use t.Errorf() to indicate that the test failed and log that
	if hd != expected {
		t.Errorf("expected %q: got %q", expected, hd)
	}
}
