package jsonp

//
// Tests that the parsing works correctly and you get back resonable values
//
// Copyright (C) Philip Schlump, 2015
// License LICENSE.apache.txt or LICENSE.mit.txt
//
//

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

const db_test = false

var TestCases = []struct {
	paths        []string
	calls        string
	input        string
	expectedBody string
}{
	{
		[]string{"/abc", "/def"},
		"/abc?callback=xyz",
		`{"ok":123}`,
		`xyz({"ok":123});`,
	},
	{
		[]string{"/abc", "/def"},
		"/ghi?callback=xyz",
		`{"ok":123}`,
		`{"ok":123}`,
	},
}

func TestIPFilter(t *testing.T) {
	for _, tc := range TestCases {

		aaa := JsonPHandlerType{
			Next:  nextFunc(tc.input),
			Paths: tc.paths,
		}

		req, err := http.NewRequest("GET", tc.calls, nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}
		req.RequestURI = tc.calls

		rec := httptest.NewRecorder()

		// Make the call to the server
		status, err := aaa.ServeHTTP(rec, req)
		if err != nil {
			t.Fatalf("Responded with error: %v, TestCase: %+v\n", err, tc)
		}
		if status != 200 {
			t.Fatalf("Responded with invalid status: %v, TestCase: %+v\n", err, tc)
		}

		resultBody := rec.Body.String()
		if db_test {
			fmt.Printf("body >%s<-\n", resultBody)
		}
		if resultBody != tc.expectedBody {
			t.Fatalf("Expected Body: '%s', Got: '%s' TestCase: %+v\n", tc.expectedBody, rec.Body.String(), tc)
		}
	}
}

func nextFunc(ss string) httpserver.Handler {
	return httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		if _, err := w.Write([]byte(ss)); err != nil {
			return 500, err
		}
		return 200, nil
	})
}
