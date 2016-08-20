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

	"github.com/mholt/caddy"
	"github.com/mholt/caddy-old1/middleware"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/pschlump/caddy-jsonp/bufferhtml"
)

const db1 = false // back to debuging this

func init() {
	caddy.RegisterPlugin("jsonp", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

// Setup configures a new Jsonp middleware instance.
func Setup(c *caddy.Controller) error {
	paths, err := jsonpParse(c)
	if err != nil {
		return err
	}

	if db1 {
		fmt.Printf("Setup called, paths=%s\n", paths)
	}

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return JsonPHandlerType{
			Next:  next,
			Paths: paths,
		}
	})

	return nil
}

func jsonpParse(c *caddy.Controller) ([]string, error) {
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
	Next  httpserver.Handler
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
				// if there is a "callback" argument then use that to format the JSONp response.
				if db1 {
					fmt.Printf("D Inside, status=%d\n", status)
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
					fmt.Printf("E - error occured, %s\n", err)
				}
				log.Printf("Error (%s) status %d - from JSONP directive\n", err, status)
			}
			return status, err
		}
	}
	return jph.Next.ServeHTTP(www, req)
}

/* vim: set noai ts=4 sw=4: */
