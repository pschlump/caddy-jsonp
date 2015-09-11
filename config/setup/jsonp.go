// Package jsonp impements JSONp middleware
//
// Some documentation on this in ../../middleware/jsonp/README.md
//
// Copyright (C) Philip Schlump, 2015
// License LICENSE.apache.txt or LICENSE.mit.txt
//

package jsonp

import (
	"net/http"
	"net/url"

	"github.com/mholt/caddy/middleware"
	"github.com/mholt/caddy/middleware/jsonp/bufferhtml"
)

type JsonPHandlerType struct {
	Paths                []string
	Next                 middleware.Handler
	ResponseBodyRecorder *bufferhtml.BufferHTML
}

func (jph JsonPHandlerType) ServeHTTP(www http.ResponseWriter, req *http.Request) (int, error) {
	for _, prefix := range jph.Paths {
		if middleware.Path(req.URL.Path).Matches(prefix) {
			ResponseBodyRecorder := bufferhtml.NewBufferHTML()
			status, err := jph.Next.ServeHTTP(ResponseBodyRecorder, req)
			if status == 200 {
				// if there is a "callback" argument then use that to format the JSONp response.
				Prefix := ""
				Postfix := ""
				u, err := url.ParseRequestURI(req.RequestURI)
				if err != nil {
					ResponseBodyRecorder.FlushAtEnd(www, "/*1*/", "")
					return status, nil
				}
				m, err := url.ParseQuery(u.RawQuery)
				if err != nil {
					ResponseBodyRecorder.FlushAtEnd(www, "/*2*/", "")
					return status, nil
				}
				callback := m.Get("callback")
				if callback != "" {
					www.Header().Set("Content-Type", "application/javascript")
					Prefix = callback + "("
					Postfix = ");"
				}
				ResponseBodyRecorder.FlushAtEnd(www, Prefix, Postfix)
			} else {
				ResponseBodyRecorder.FlushAtEnd(www, "/*3*/", "")
			}
			return status, err
		}
	}
	return jph.Next.ServeHTTP(www, req)
}
