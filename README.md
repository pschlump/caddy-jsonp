# caddy-jsonp -- A JSONP Implementation for the Caddy Webserver in Go (GoLang)

 [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/pschlump/Go-FTL/master/LICENSE)

This is both an implementation of JSONP for Caddy and a description of the process
of building middleware for Caddy.  (Please note the code has been updated for the 0.9
version of Caddy - but the documentation in this file may be out of date.  I will
get to it tomorrow morning.)

[Caddy](https://github.com/mholt/caddy "Caddy Web Server") is a fast web server and
proxy that supports middleware.  The middleware is somewhat different from the
regular middleware in a Go webserver.


# Test

Run caddy with the ./example/Caddyfile.

You can now fetch www/status.json in a JSONP format with

``` sh

	$ wget --no-check-certificate -o aa.out -O bb.out \
		'https://localhost/api/status?callback=callback11212121'

```

Will return the results of /api/status in the return value of

``` javascript

	callback11212121({"status":"ok"});

```

This is compatible with using JSONP in jQuery and AngularJS.


## Things to note

This directive is not compatible with streaming calls.  This should not be a
problem because streaming is not compatible with JSONP to start off with.
JSONP requires that the complete set of data be in the parameters to the
callback function.  Streaming only returns the data in chunks.

The code internally buffers the entire response in memory before sending
it back.  Very large responses should be avoided.

## Example Call

Example call using jQuery

``` jvascript

		$('#email_submit').click(function(){
			var email_addr = $("#email_address").val();
			var email_subject = $("#email_subject").val();
			var email_body = $("#email_body").val();
			$.ajax({
				url: "http://www.2c-why.com/testapi/send-email"
				, jsonp: "callback" // The name of the callback parameter
				, dataType: "jsonp" // Tell jQuery we're expecting JSONP
				, data: { "email_addr": email_addr, "subject": email_subject, "body": email_body, "key": global_email_key }
				, success: function( resp ) {
					if ( resp.status === "success" ) {
						$("#success_msg").show();	
					} else {
						$('#message_error').show();
					}
				}
				, error: function( resp ) {
					$('#message_error').show();
				}
			});
			return false;
		});

```


## To generate certificates 

To generate your own self signed certificates you can:

``` base

	$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
	$ mv *.pem ./key

```

## Author

Philip Schlump

[Authors Website, 2C Why, LLC](http://www.2c-why.com/ "2C Why LLC")


