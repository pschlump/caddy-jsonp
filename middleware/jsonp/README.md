# JSONp middleware and directive

This implements JSONP for Caddy - a full featured fast web server.

## Usage

In the `Caddyfile` add in

```

	jsonp /some/url.json
	jsonp /some/other/url/

```

This will turn on the ability to fetch these URLs with a GET in the JSONP protocol.
It will be returned in JSONP if an get parameter of `callback` is supplied.  For
example, given that a call to /api/status returns JSON, lets say `{"status":"ok"}`

``` sh

	$ wget -o aa.out -O bb.out 'http://localhost/api/status?callback=callback11212121'

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

## Author

Philip Schlump






	

