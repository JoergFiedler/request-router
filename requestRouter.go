package rero

import (
	"net/http"
)

// Handler function definition
type Handler func(RequestContext)

// RequestContext passed to handlers used to pass values down the
// handler chain and respond to the client
type RequestContext interface {
	// get original http request
	GetRequest() *http.Request
	// set context variable value
	SetVar(name string, value interface{})
	// get requests contexts value for variable name
	GetVar(name string) interface{}
	// break handler chain and respond with status and body
	Respond(status int, body string)
}

// RequestRouter used to add handlers to routes
type RequestRouter interface {
	AddGlobalHandler(string, Handler)
	AddPathHandler(string, string, Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// CreateRouter creates a router instance
func CreateRouter() RequestRouter {
	return &routeConfig{
		pathTree:       createPathTree(),
		globalHandlers: make([]Handler, 0),
	}
}

type handlerContext struct {
	request *http.Request
	vars    map[string]interface{}
	closed  bool
	status  int
	body    string
}

type routeConfig struct {
	globalHandlers []Handler
	pathTree       pathTree
}

func (config *routeConfig) AddGlobalHandler(
	method string,
	handler Handler,
) {
	config.globalHandlers = append(config.globalHandlers, handler)
}

func (config *routeConfig) AddPathHandler(
	method string,
	path string,
	handler Handler,
) {
	config.pathTree.addPathHandler(method, path, handler)
}

func (config *routeConfig) ServeHTTP(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	pathContext := config.pathTree.getPathContext(
		request.Method,
		request.URL.Path,
	)

	if len(pathContext.handlers) == 0 {
		responseWriter.WriteHeader(404)
		responseWriter.Write([]byte("Not Found"))
		return
	}
	handlers := append(config.globalHandlers, pathContext.handlers...)
	vars := make(map[string]interface{}, 0)
	for key, value := range pathContext.pathVariables {
		vars[key] = value
	}

	requestContext := &handlerContext{request: request, vars: vars }
	for _, handler := range handlers {
		handler(requestContext)
		if requestContext.closed {
			break
		}
	}

	responseWriter.WriteHeader(requestContext.status)
	responseWriter.Write([]byte(requestContext.body))
}

func (handlerContext *handlerContext) Respond(status int, body string) {
	handlerContext.closed = true
	handlerContext.status = status
	handlerContext.body = body
}

func (handlerContext *handlerContext) GetRequest() *http.Request {
	return handlerContext.request
}

func (handlerContext *handlerContext) SetVar(name string, value interface{}) {
	handlerContext.vars[name] = value
}

func (handlerContext *handlerContext) GetVar(name string) interface{} {
	return handlerContext.vars[name]
}
