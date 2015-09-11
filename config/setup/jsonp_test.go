package setup

// Tests that the parsing works correctly and you get back resonable values
// Modified by Philip Schlump from ./internal_test.go
// Some documentation on this in ../../middleware/jsonp/README.md

import (
	"testing"

	"github.com/mholt/caddy/middleware/jsonp"
)

func TestJsonp(t *testing.T) {
	c := NewTestController(`jsonp /api/status`)

	mid, err := Jsonp(c)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	if mid == nil {
		t.Fatal("Expected middleware, was nil instead")
	}

	handler := mid(EmptyNext)
	myHandler, ok := handler.(jsonp.JsonPHandlerType)

	if !ok {
		t.Fatalf("Expected handler to be type JsonPHandlerType, got: %#v", handler)
	}

	if myHandler.Paths[0] != "/api/status" {
		t.Errorf("Expected /api/status in the list of Jsonp Paths")
	}

	if !SameNext(myHandler.Next, EmptyNext) {
		t.Error("'Next' field of handler was not set properly")
	}

}

func TestJsonpParse(t *testing.T) {
	tests := []struct {
		inputJsonpPaths    string
		shouldErr          bool
		expectedJsonpPaths []string
	}{
		{`jsonp /api/status`, false, []string{"/api/status"}},

		{`jsonp /api/status
		  jsonp /api/email`, false, []string{"/api/status", "/api/email"}},
	}
	for i, test := range tests {
		c := NewTestController(test.inputJsonpPaths)
		actualJsonpPaths, err := jsonpParse(c)

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if len(actualJsonpPaths) != len(test.expectedJsonpPaths) {
			t.Fatalf("Test %d expected %d JsonpPaths, but got %d",
				i, len(test.expectedJsonpPaths), len(actualJsonpPaths))
		}
		for j, actualJsonpPath := range actualJsonpPaths {
			if actualJsonpPath != test.expectedJsonpPaths[j] {
				t.Fatalf("Test %d expected %dth Jsonp Path to be  %s  , but got %s",
					i, j, test.expectedJsonpPaths[j], actualJsonpPath)
			}
		}
	}

}
