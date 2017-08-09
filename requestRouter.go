package rero

import (
	"net/http"
)

type Handler func(RequestContext)

type HandlerContext struct {
	request *http.Request
	vars    map[string]interface{}
	closed  bool
	status  int
	body    string
}

type RequestContext interface {
	GetRequest() *http.Request
	GetVars() map[string]interface{}
	GetVar(name string) interface{}
	Respond(int, string)
}

type RouteConfig struct {
	globalHandlers []Handler
	pathTree       PathTree
}

type RequestRouter interface {
	AddGlobalHandler(string, Handler)
	AddPathHandler(string, string, Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func CreateRouter() RequestRouter {
	return &RouteConfig{
		pathTree:       CreatePathTree(),
		globalHandlers: make([]Handler, 0),
	}
}

func (config *RouteConfig) AddGlobalHandler(
	method string,
	handler Handler,
) {
	config.globalHandlers = append(config.globalHandlers, handler)
}

func (config *RouteConfig) AddPathHandler(
	method string,
	path string,
	handler Handler,
) {
	config.pathTree.AddPathHandler(method, path, handler)
}

func (config *RouteConfig) ServeHTTP(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	pathContext := config.pathTree.GetPathContext(
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

	requestContext := &HandlerContext{request: request, vars: vars }
	for _, handler := range handlers {
		handler(requestContext)
		if requestContext.closed {
			break
		}
	}

	responseWriter.WriteHeader(requestContext.status)
	responseWriter.Write([]byte(requestContext.body))
}

func (handlerContext *HandlerContext) Respond(status int, body string) {
	handlerContext.closed = true
	handlerContext.status = status
	handlerContext.body = body
}

func (handlerContext *HandlerContext) GetRequest() *http.Request {
	return handlerContext.request
}

func (handlerContext *HandlerContext) GetVars() map[string]interface{} {
	return handlerContext.vars
}

func (handlerContext *HandlerContext) GetVar(name string) interface{} {
	return handlerContext.vars[name]
}
