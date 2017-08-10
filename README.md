# Request-Router

[![Build Status](https://travis-ci.org/JoergFiedler/request-router.svg?branch=master)](https://travis-ci.org/JoergFiedler/request-router)
[![Coverage Status](https://coveralls.io/repos/github/JoergFiedler/request-router/badge.svg)](https://coveralls.io/github/JoergFiedler/request-router)

Simple request router. This is actually my training project to get acquainted
with go language. Feedback and suggestions are highly appreciated.

Features:

- register multiple handlers for a single route
- match requests to routes based on path and method
- match path variables
- register global handlers  

Next steps:

- improve interface for registering handlers

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/JoergFiedler/rero
```

## Quick Start

```go
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
```

Done.

## Match Requests

```go
// matches all GET requests with path / 
router.AddPathHandler("GET", "/", handler)
// matches all HEAD requests with path /hello and /hello/
router.AddPathHandler("HEAD", "/hello/", handler)
// matches all POST requests with path /hello/world and /hello/world/
router.AddPathHandler("POST", "/hello/world", handler)
```

## Path Variables

Path variables are supported. The following example will respond with `Hello, Max.`,
if requested with `curl http://localhost:4275/hello/Max`.

```go
handler := func(ctx rero.RequestContext) {
  ctx.Respond(200, fmt.Sprintf("Hello, %s.", ctx.GetVar("name")))
}

router.AddPathHandler("GET", "/hello/:name:", handler)
```

Multiple variables within a path are supported as well.

```go
router.AddPathHandler("GET", "/hello/:name:/:age:", handler)
```

It is not possible to have both variable and non-variable path segments. The
router will panic when you try to register the second handler. No way
to decided which handler should be called.

```go
router.AddPathHandler("GET", "/hello/:name:", handler)
router.AddPathHandler("GET", "/hello/world", handler) // does not work, panic
router.AddPathHandler("GET", "/hi/world", handler) // does work
```

## Multiple Handlers

The following snippet is also valid. The handlers will be called in order they
were added.

```go
router.AddPathHandler("GET", "/hello/world", handler1)
router.AddPathHandler("GET", "/hello/world", handler2)
```

Handlers can communicate via `RequestContext.GetVars()`. A handler can add
objects to the var map, which then can be used by subsequent handlers. If
a handlers calls `RequestContext.Respond()` no other handler will be called
anymore.
   
## Global Handlers

Global handlers will always be executed before all other handlers. They also
can add variables to `RequestContext`. This can be used to validate auth tokens.
If the global handler calls `RequestContext.Respond()` no other handler will
be executed anymore.

```go
router.AddGlobalHandler("GET", authHandler)
```