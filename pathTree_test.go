package rero

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tree pathTree

var _ = Describe("pathTree", func() {
	Context("createPathTree", func() {
		BeforeEach(func() {
			tree = createPathTree()
		})

		It("returns an instance", func() {
			Expect(tree).NotTo(BeNil())
		})
	})

	Context("getPathContext", func() {
		handler := func(any RequestContext) {}

		BeforeEach(func() {
			tree = createPathTree()
		})

		It("returns list of existing handlers", func() {
			tree.addPathHandler("GET", "/a", handler).addPathHandler("GET", "/a", handler)

			Expect(len(tree.getPathContext("GET", "/a").handlers)).To(Equal(2))
		})

		It("ignores trailing slash", func() {
			tree.addPathHandler("GET", "/", handler).addPathHandler("GET", "", handler)

			Expect(len(tree.getPathContext("GET", "/").handlers)).To(Equal(2))
		})

		It("compares complete path", func() {
			tree.addPathHandler("GET", "/a/b/c", handler).addPathHandler("GET", "/z/b/c", handler)

			Expect(len(tree.getPathContext("GET", "/a/b/c").handlers)).To(Equal(1))
			Expect(len(tree.getPathContext("GET", "/z/b/c").handlers)).To(Equal(1))
		})

		It("skips intermediate path segments", func() {
			tree.addPathHandler("GET", "/a/b/c", handler)

			Expect(len(tree.getPathContext("GET", "/").handlers)).To(Equal(0))
			Expect(len(tree.getPathContext("GET", "/a").handlers)).To(Equal(0))
			Expect(len(tree.getPathContext("GET", "/a/b").handlers)).To(Equal(0))
		})

		It("respects method selector", func() {
			tree.addPathHandler("PUT", "/a/b/c", handler)

			Expect(len(tree.getPathContext("GET", "/a/b/c").handlers)).To(Equal(0))
			Expect(len(tree.getPathContext("PUT", "/a/b/c").handlers)).To(Equal(1))
		})

		It("matches longest path only", func() {
			tree.addPathHandler("GET", "/a", handler).
				addPathHandler("GET", "/a/b", handler)

			Expect(len(tree.getPathContext("GET", "/a/b").handlers)).To(Equal(1))
		})

		It("returns an empty array for unregistered paths", func() {
			Expect(len(tree.getPathContext("GET", "").handlers)).To(Equal(0))
			Expect(len(tree.getPathContext("GET", "/a").handlers)).To(Equal(0))
			Expect(len(tree.getPathContext("GET", "/z/yy").handlers)).To(Equal(0))
		})

		It("replaces path variables", func() {
			tree.addPathHandler("GET", "/a/:var1:/c", handler)
			tree.addPathHandler("GET", "/a/:var1:/a/d/:var2:/f", handler)

			Expect(tree.getPathContext("GET", "/a/alice/c").pathVariables).
				To(Equal(map[string]string{"var1": "alice"}))
			Expect(tree.getPathContext("GET", "/a/john/a/d/murdoch/f").pathVariables).
				To(Equal(map[string]string{"var1": "john", "var2": "murdoch"}))
		})

		It("panic if non-var segment clashes with var segment", func() {
			tree.addPathHandler("GET", "/a/:var1:/c", handler)

			addNonUniformPath := func() {
				tree.addPathHandler("GET", "/a/b/c", handler)
			}
			Expect(addNonUniformPath).To(Panic())
		})

		It("panic if var segment clashes with non-var segment", func() {
			tree.addPathHandler("GET", "/a/b/c", handler)

			addNonUniformPath := func() {
				tree.addPathHandler("GET", "/a/:var1:/c", handler)
			}
			Expect(addNonUniformPath).To(Panic())
		})
	})
})
