package router

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var pathTree PathTree

var _ = Describe("PathTree", func() {
	Context("CreatePathTree", func() {
		BeforeEach(func() {
			pathTree = CreatePathTree()
		})

		It("returns an instance", func() {
			Expect(pathTree).NotTo(BeNil())
		})
	})

	Context("GetPathContext", func() {
		handler := func(any RequestContext) {}

		BeforeEach(func() {
			pathTree = CreatePathTree()
		})

		It("returns list of existing handlers", func() {
			pathTree.AddPathHandler("GET", "/a", handler).AddPathHandler("GET", "/a", handler)

			Expect(len(pathTree.GetPathContext("GET", "/a").handlers)).To(Equal(2))
		})

		It("ignores trailing slash", func() {
			pathTree.AddPathHandler("GET", "/", handler).AddPathHandler("GET", "", handler)

			Expect(len(pathTree.GetPathContext("GET", "/").handlers)).To(Equal(2))
		})

		It("compares complete path", func() {
			pathTree.AddPathHandler("GET", "/a/b/c", handler).AddPathHandler("GET", "/z/b/c", handler)

			Expect(len(pathTree.GetPathContext("GET", "/a/b/c").handlers)).To(Equal(1))
			Expect(len(pathTree.GetPathContext("GET", "/z/b/c").handlers)).To(Equal(1))
		})

		It("skips intermediate path segments", func() {
			pathTree.AddPathHandler("GET", "/a/b/c", handler)

			Expect(len(pathTree.GetPathContext("GET", "/").handlers)).To(Equal(0))
			Expect(len(pathTree.GetPathContext("GET", "/a").handlers)).To(Equal(0))
			Expect(len(pathTree.GetPathContext("GET", "/a/b").handlers)).To(Equal(0))
		})

		It("respects method selector", func() {
			pathTree.AddPathHandler("PUT", "/a/b/c", handler)

			Expect(len(pathTree.GetPathContext("GET", "/a/b/c").handlers)).To(Equal(0))
			Expect(len(pathTree.GetPathContext("PUT", "/a/b/c").handlers)).To(Equal(1))
		})

		It("matches longest path only", func() {
			pathTree.AddPathHandler("GET", "/a", handler).
				AddPathHandler("GET", "/a/b", handler)

			Expect(len(pathTree.GetPathContext("GET", "/a/b").handlers)).To(Equal(1))
		})

		It("returns an empty array for unregistered paths", func() {
			Expect(len(pathTree.GetPathContext("GET", "").handlers)).To(Equal(0))
			Expect(len(pathTree.GetPathContext("GET", "/a").handlers)).To(Equal(0))
			Expect(len(pathTree.GetPathContext("GET", "/z/yy").handlers)).To(Equal(0))
		})

		It("replaces path variables", func() {
			pathTree.AddPathHandler("GET", "/a/:var1:/c", handler)
			pathTree.AddPathHandler("GET", "/a/:var1:/a/d/:var2:/f", handler)

			Expect(pathTree.GetPathContext("GET", "/a/alice/c").pathVariables).
				To(Equal(map[string]string{"var1": "alice"}))
			Expect(pathTree.GetPathContext("GET", "/a/john/a/d/murdoch/f").pathVariables).
				To(Equal(map[string]string{"var1": "john", "var2": "murdoch"}))
		})

		It("panic if non-var segment clashes with var segment", func() {
			pathTree.AddPathHandler("GET", "/a/:var1:/c", handler)

			addNonUniformPath := func() {
				pathTree.AddPathHandler("GET", "/a/b/c", handler)
			}
			Expect(addNonUniformPath).To(Panic())
		})

		It("panic if var segment clashes with non-var segment", func() {
			pathTree.AddPathHandler("GET", "/a/b/c", handler)

			addNonUniformPath := func() {
				pathTree.AddPathHandler("GET", "/a/:var1:/c", handler)
			}
			Expect(addNonUniformPath).To(Panic())
		})
	})
})
