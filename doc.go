// See LICENSE file.

/*
Package rero implements a simple request router and dispatcher.

Features:

    * register multiple handlers for a single route
    * match requests to routes based on path and method
    * match path variables
    * register global handlers

Example:

    package main

    import (
      "net/http"

      "github.com/JoergFiedler/rero"
    )

    func main() {

      // Create router instance
      router := rero.CreateRouter()


      // This does not do much so lets register a simple handler. A handler
      // looks a bit like this.
      handler := func(ctx rero.RequestContext) {
        ctx.Respond(200, "Hello there.")
      }

      // Lets handle all request to `/hello` path.
      router.AddPathHandler("GET", "/hello", handler)

      // Marry the router with `http` package.
      http.ListenAndServe(":4275", router)
    }

Requests are matched by method and path. Trailing slashes are ignored.

    // matches all GET requests with path /
    router.AddPathHandler("GET", "/", handler)
    // matches all HEAD requests with path /hello and /hello/
    router.AddPathHandler("HEAD", "/hello/", handler)
    // matches all POST requests with path /hello/world and /hello/world/
    router.AddPathHandler("POST", "/hello/world", handler)

Path variables are supported as well. Use as many variables within in single path as you like.

    handler := func(ctx rero.RequestContext) {
      ctx.Respond(200, fmt.Sprintf("Hello, %s.", ctx.GetVar("name")))
    }
    router.AddPathHandler("GET", "/hello/:name:", handler)

Configuring routes that have variable and non-variable path segments is not supported. The router
will panic when you try to do this. It would be impossible to decide which handler to call.

    router.AddPathHandler("GET", "/hello/:name:", handler)
    router.AddPathHandler("GET", "/hello/world", handler) // does not work, panic
    router.AddPathHandler("GET", "/hi/world", handler) // does work

You may register multiple handler per route. All handlers will be called in order they were registered. The
chain breaks when a handler calls rero.RequestContext.Respond().

    router.AddPathHandler("GET", "/hello/world", handler1)
    router.AddPathHandler("GET", "/hello/world", handler2)

To register global handlers, which will be executed before route specific handlers. Use
rero.RequestContext.SetVar() and rero.RequestContext.GetVar() to communicated between
handlers.

    authHandler := func(ctx rero.RequestContext) {
      ctx.SetVar("username", "john")
    }
    router.AddGlobalHandler("GET", authHandler)
*/
package rero
