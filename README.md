# caddy-jsonp -- A JSONP Implementation for the Caddy Webserver in Go (GoLang)

This is both an implementation of JSONP for Caddy and a description of the process
of building middleware for Caddy.

[Caddy](https://github.com/mholt/caddy "Caddy Web Server") is a fast web server and
proxy that supports middleware.  The middleware is somewhat different from the
regular middleware in a Go webserver.

The author of the software recommends the use of caddydev a tool for helping build
middleware.  I was unable to get this to work (the tool is still under development
at this time).  For that reason I have gone ahead with building this middleware
without the additional tool.




