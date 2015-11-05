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
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mholt/caddy/caddy/setup" // Old: version 0.7.6 "github.com/mholt/caddy/config/setup"
	"github.com/mholt/caddy/middleware"  //  Old: version 0.7.6 "github.com/mholt/caddy/middleware"
	"github.com/pschlump/caddy-jsonp/bufferhtml"
)

const db1 = false

// Jsonp configures a new Jsonp middleware instance.
func Setup(c *setup.Controller) (middleware.Middleware, error) {
	paths, err := jsonpParse(c)
	if err != nil {
		return nil, err
	}

	if db1 {
		fmt.Printf("Setup called, paths=%s\n", paths)
	}

	return func(next middleware.Handler) middleware.Handler {
		return JsonPHandlerType{
			Next:  next,
			Paths: paths,
		}
	}, nil
}

func jsonpParse(c *setup.Controller) ([]string, error) {
	var paths []string

	if db1 {
		fmt.Printf("jsonpParse called\n")
	}

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

	if db1 {
		fmt.Printf("JsonPHandlerType.ServeHTTP called\n")
	}

	for _, prefix := range jph.Paths {
		if db1 {
			fmt.Printf("Path Matches\n")
		}
		if middleware.Path(req.URL.Path).Matches(prefix) {
			if db1 {
				fmt.Printf("A\n")
			}
			ResponseBodyRecorder := bufferhtml.NewBufferHTML()
			if db1 {
				fmt.Printf("B\n")
			}
			status, err := jph.Next.ServeHTTP(ResponseBodyRecorder, req)
			if db1 {
				fmt.Printf("C status=%d err=%s\n", status, err)
			}
			if status == 200 || status == 0 {
				if db1 {
					fmt.Printf("D\n")
				}
				// if there is a "callback" argument then use that to format the JSONp response.
				if db1 {
					fmt.Printf("Inside, status=%d\n", status)
				}
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
					if db1 {
						fmt.Printf("Callback = -[%s]-\n", callback)
					}
					www.Header().Set("Content-Type", "application/javascript")
					Prefix = callback + "("
					Postfix = ");"
				}
				ResponseBodyRecorder.FlushAtEnd(www, Prefix, Postfix)
			} else {
				if db1 {
					fmt.Printf("E - error occured\n")
				}
				// ResponseBodyRecorder.FlushAtEnd(www, "", "")
				log.Printf("Error (%s) status %d - from JSONP directive\n", err, status)
			}
			return status, err
		}
	}
	return jph.Next.ServeHTTP(www, req)
}
