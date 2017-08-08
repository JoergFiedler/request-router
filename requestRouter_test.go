package router

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockResponseWriter struct {
	header   map[string][]string
	status   int
	response []byte
}

var requestRouter RequestRouter

var _ = Describe("Router", func() {
	Context("CreatePathTree", func() {
		BeforeEach(func() {
			requestRouter = CreateRouter()
		})

		It("returns an instance", func() {
			Expect(requestRouter).NotTo(BeNil())
		})
	})

	Context("AddPathHandler", func() {
		BeforeEach(func() {
			requestRouter = CreateRouter()
		})

		It("adds a handler", func() {
			request := &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
			responseWriter := createMockResponseWriter()
			handler := func(rc RequestContext) {
				Expect(rc.GetRequest()).To(Equal(request))
			}

			requestRouter.AddPathHandler("GET", "/", handler)
			requestRouter.ServeHTTP(responseWriter, request)
		})
	})

	Context("ServeHTTP", func() {
		var request *http.Request
		var responseWriter *MockResponseWriter

		BeforeEach(func() {
			requestRouter = CreateRouter()
			request = &http.Request{Method: "GET", URL: &url.URL{Path: "/"}}
			responseWriter = createMockResponseWriter()
		})

		It("returns 404 for unhandled paths", func() {
			requestRouter.ServeHTTP(responseWriter, request)

			Expect(responseWriter.status).To(Equal(404))
			Expect(responseWriter.response).To(Equal([]byte("Not Found")))
		})

		It("responds with result from handler", func() {
			status := 301
			response := "Temporarily Redirect"
			handler := func(rc RequestContext) {
				rc.Respond(status, response)
			}

			requestRouter.AddPathHandler("GET", "/", handler)
			requestRouter.ServeHTTP(responseWriter, request)

			Expect(responseWriter.status).To(Equal(status))
			Expect(responseWriter.response).To(Equal([]byte(response)))
		})

		It("calls all registered handlers in perfect order", func() {
			calls := []string{}
			handler1 := func(rc RequestContext) {
				calls = append(calls, "handler1")
			}
			handler2 := func(rc RequestContext) {
				calls = append(calls, "handler2")
			}

			requestRouter.AddPathHandler("GET", "/", handler1)
			requestRouter.AddPathHandler("GET", "/", handler2)
			requestRouter.ServeHTTP(responseWriter, request)

			Expect(calls).To(Equal([]string{"handler1", "handler2"}))
		})

		It("passes on context vars from handler to handler", func() {
			handler1 := func(rc RequestContext) {
				rc.GetVars()["john"] = "doe"
			}
			handler2 := func(rc RequestContext) {
				Expect(rc.GetVars()["john"]).To(Equal("doe"))
			}

			requestRouter.AddPathHandler("GET", "/", handler1)
			requestRouter.AddPathHandler("GET", "/", handler2)

			requestRouter.ServeHTTP(responseWriter, request)
		})

		It("calls global handlers first", func() {
			callOrder := make([]string, 0)
			globalHandler := func(rc RequestContext) {
				callOrder = append(callOrder, "global")
			}
			pathHandler := func(rc RequestContext) {
				callOrder = append(callOrder, "path")
			}

			requestRouter.AddPathHandler("GET", "/", pathHandler)
			requestRouter.AddGlobalHandler("GET", globalHandler)

			requestRouter.ServeHTTP(responseWriter, request)

			Expect(callOrder).To(Equal([]string{"global", "path"}))
		})

		It("add path variables to context vars", func() {
			handler := func(rc RequestContext) {
				Expect(rc.GetVars()["entity"]).To(Equal("elephant"))
				Expect(rc.GetVar("entity")).To(Equal("elephant"))
			}

			request.URL.Path = "/elephant/"
			requestRouter.AddPathHandler("GET", "/:entity:/", handler)

			requestRouter.ServeHTTP(responseWriter, request)
		})
	})
})

func createMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{}
}

func (mock *MockResponseWriter) Header() http.Header {
	return mock.header
}

func (mock *MockResponseWriter) WriteHeader(status int) {
	mock.status = status
}

func (mock *MockResponseWriter) Write(response []byte) (int, error) {
	mock.response = response

	return 0, nil
}
