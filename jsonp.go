// Package jsonp impements JSONp middleware
//
// Copyright (C) Philip Schlump, 2015
// License LICENSE.apache.txt or LICENSE.mit.txt
//

//
// Directive
//
// jsonp /some/url
// jsonp /some/other-url
//

package jsonp

import (
	"net/http"
	"net/url"

	"github.com/mholt/caddy/config/setup"
	"github.com/mholt/caddy/middleware"
	"github.com/pschlump/caddy-jsonp/bufferhtml"
)

// Jsonp configures a new Jsonp middleware instance.
func Setup(c *setup.Controller) (middleware.Middleware, error) {
	paths, err := jsonpParse(c)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Setup called, paths=%s\n", paths)

	return func(next middleware.Handler) middleware.Handler {
		return JsonPHandlerType{
			Next:  next,
			Paths: paths,
		}
	}, nil
}

func jsonpParse(c *setup.Controller) ([]string, error) {
	var paths []string

	// fmt.Printf("jsonpParse called\n")

	for c.Next() {
		if !c.NextArg() {
			return paths, c.ArgErr()
		}
		paths = append(paths, c.Val())
	}

	return paths, nil
}

type JsonPHandlerType struct {
	Paths []string
	Next  middleware.Handler
}

func (jph JsonPHandlerType) ServeHTTP(www http.ResponseWriter, req *http.Request) (int, error) {

	// fmt.Printf("JsonPHandlerType.ServeHTTP called\n")

	for _, prefix := range jph.Paths {
		if middleware.Path(req.URL.Path).Matches(prefix) {
			ResponseBodyRecorder := bufferhtml.NewBufferHTML()
			status, err := jph.Next.ServeHTTP(ResponseBodyRecorder, req)
			if status == 200 || status == 0 {
				// if there is a "callback" argument then use that to format the JSONp response.
				Prefix := ""
				Postfix := ""
				u, err := url.ParseRequestURI(req.RequestURI)
				if err != nil {
					ResponseBodyRecorder.FlushAtEnd(www, "", "")
					return status, nil
				}
				m, err := url.ParseQuery(u.RawQuery)
				if err != nil {
					ResponseBodyRecorder.FlushAtEnd(www, "", "")
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
				ResponseBodyRecorder.FlushAtEnd(www, "", "")
			}
			return status, err
		}
	}
	return jph.Next.ServeHTTP(www, req)
}
