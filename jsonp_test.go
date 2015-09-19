package jsonp

//
// Tests that the parsing works correctly and you get back resonable values
// Some documentation on this in ../../middleware/jsonp/README.md
//
// Copyright (C) Philip Schlump, 2015
// License LICENSE.apache.txt or LICENSE.mit.txt
//

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mholt/caddy/middleware"
)

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
}

func TestIPFilter(t *testing.T) {
	for _, tc := range TestCases {
		aaa := JsonPHandlerType{
			Next: middleware.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
				return http.StatusOK, nil
			}),
			Paths: tc.paths,
		}

		req, err := http.NewRequest("GET", tc.calls, nil)
		if err != nil {
			t.Fatalf("Could not create HTTP request: %v", err)
		}

		rec := httptest.NewRecorder()

		status, err := aaa.ServeHTTP(rec, req)
		if err != nil {
			t.Fatalf("Responded with error: %v, TestCase: %v\n", err, tc)
		}
		if status != 200 {
			t.Fatalf("Responded with invalid status: %v, TestCase: %v\n", err, tc)
		}

		if rec.Body.String() != tc.expectedBody {
			t.Fatalf("Expected Body: '%s', Got: '%s' TestCase: %v\n", tc.expectedBody, rec.Body.String(), tc)
		}
	}
}
