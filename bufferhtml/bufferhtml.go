//
// Copyright (C) Philip Schlump, 2015
// License LICENSE.apache.txt or LICENSE.mit.txt
//

package bufferhtml

import (
	"bytes"
	"fmt"
	"net/http"
)

// ------------------------------------------------------------------------------------------------------------
// Implement a compatible http.ResponseWriter that saves all the data until the end.
// Good side: you don't need to finish your headers before your data.
// Good side: you can post-process the data/status.
// Good side: length header can be manipulated after the data is generated.
// Bad side: This won't work with a streaming data interface at all.
// Bad side: Also it's all buffered in memory.
type BufferHTML struct {
	bytes.Buffer             // The body of the response
	StatusCode   int         // StatusCode like 200, 404 etc.
	Headers      http.Header // All the headers that will be writen when done
}

// Return a new buffered http.ResponseWriter
func NewBufferHTML() *BufferHTML {
	return &BufferHTML{
		Headers:    make(http.Header),
		StatusCode: 0,
	}
}

// Return the headers - Required to make the interface work
func (b *BufferHTML) Header() http.Header {
	return b.Headers
}

// Implement http.ResponseWriter WriteHeader to just buffer the Status
func (b *BufferHTML) WriteHeader(StatusCode int) {
	b.StatusCode = StatusCode
}

func (b *BufferHTML) FlushAtEnd(w http.ResponseWriter, Prefix string, Postfix string) (n int, err error) {
	h := w.Header()
	s := b.Bytes()
	if len(b.Headers) > 0 {
		for key, val := range b.Headers {
			h[key] = val
		}
	}
	// ------------------------------------------- prefix / postfix --------------------------------
	s = []byte(Prefix + string(s) + Postfix)
	l := len(s)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", l))
	if b.StatusCode > 0 {
		w.WriteHeader(b.StatusCode)
	} else {
		w.WriteHeader(200)
	}
	n, err = w.Write(s)
	return
}
