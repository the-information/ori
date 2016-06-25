
# tutorial

This is a walkthrough of how to use ori for beginners to Go and App Engine. Note that it expects
at least a passing familiarity with Go as a language; if there are aspects of this tutorial that
confuse you, I recommend [A Tour of Go](https://tour.golang.org).

## install things

Start by installing both [Go](https://golang.org/dl/) and the [App Engine dev tools for Go](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go). If you're on a Mac with Homebrew, this is as simple as:

```bash
$ brew install go app-engine-go-64
```

## set up workspace

Then set up a $GOPATH by adding it to your profile (`.bash_profile` or `.zshrc`):
```bash
export GOPATH=/Users/goldibex/Development/go    # make sure you mkdir -p this path
```

and open a new shell so the new settings have taken effect.

Then create a new directory for your app inside your `$GOPATH`. Note that if you're developing on Github, you probably want to use the go convention for this:

```bash
$ cd $GOPATH
$ mkdir -p github.com/goldibex/somenewapi && cd github.com/goldibex/somenewapi
```

You can run `git init` here as well.

Now create a new App Engine app by instantiating an `app.yaml`. This is a file that
[App Engine uses](https://cloud.google.com/appengine/docs/go/config/appref)
to set up your application. For now, we'll just tell it to send literally every HTTP request
to our Go application:

```yaml
application: somenewapi
version: "1"
runtime: go
api_version: go1

handlers:
  - url: /.*
    script: _go_app
```

## server.go

Let's write an HTTP server in the file `server.go. This will be the core
of our application, where we tell App Engine to use kami to route requests.
We'll also set up ori's middleware handlers to intercept inbound requests, check them
for validity, and make application-wide configuration available.

```go
package app

import (
	// kami gives us URL routing by method and parameterized path,
	// as well as a convenient way to write request handlers.
	"github.com/guregu/kami"
	// ori/config provides application-wide configuration.
	"github.com/the-information/ori/config"
	// ori/rest provides content negotiation and CORS support.
	"github.com/the-information/ori/rest"
	"net/http"
)

func init() {

	// When somebody tries to GET a route that only has a POST handler,
	// respond with 405 Method Not Allowed rather than 404 Not Found.
	kami.EnableMethodNotAllowed(true)
	// Get ori to load app configuration on a per-request basis.
	kami.Use("/", config.Middleware)
	// Get ori to validate all requests as application/json encoded in
	// UTF-8.
	kami.Use("/", rest.Middleware)

	// When a request comes up as 405 Method Not Allowed, send a
	// JSON message explaining the problem.
	kami.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rest.WriteJSON(w, &rest.ErrMethodNotAllowed)
	})

	// When a request comes up as 404 Not Found because of the router,
	// send a JSON message explaining the problem.
	kami.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rest.WriteJSON(w, &rest.ErrNotFound)
	})

	// Install Kami as the default HTTP handler. App Engine
	// will take over from here.
	http.Handle("/", kami.Handler())

}
```

We can now run this server with the command `goapp get && goapp serve .`:

```bash
$ goapp serve .
INFO     2016-06-25 15:27:34,605 devappserver2.py:769] Skipping SDK update check.
INFO     2016-06-25 15:27:34,644 api_server.py:205] Starting API server at: http://localhost:37781
INFO     2016-06-25 15:27:34,739 dispatcher.py:197] Starting module "default" running at: http://localhost:8080
INFO     2016-06-25 15:27:34,740 admin_server.py:116] Starting admin server at: http://localhost:8000
```

Lovely, our server now seems to be running. Let's try accessing something:

```bash
$ curl -i localhost:8080
HTTP/1.1 404 Not Found
accept-charset: UTF-8
content-type: application/json; charset=UTF-8
accept: application/json
Cache-Control: no-cache
Expires: Fri, 01 Jan 1990 00:00:00 GMT
Content-Length: 0
Server: Development/2.0
Date: Sat, 25 Jun 2016 15:30:14 GMT
{"message":"The requested resource was not located on this server."}
```

It's working, but we haven't written any routes yet, so everything will
come up 404 Not Found.

## next up

Future tutorials will include:
- How to write route handlers,
- How to secure routes using [ori/account/auth](https://godoc.org/github.com/the-information/ori/account/auth),
- How to to store and retrieve data from App Engine Datastore,
- How to log to App Engine,
- How to conduct work outside of the request context using Task Queues.
