package router

import (
	"fmt"
	"regexp"
	"strings"
)

type PathNode struct {
	segment      string
	children     map[string]*PathNode
	pathHandlers map[string][]Handler
}

type PathTree interface {
	AddPathHandler(method string, path string, handler Handler) PathTree
	GetPathContext(method string, path string) PathContext
}

type PathContext struct {
	handlers      []Handler
	pathVariables map[string]string
}

var varRegexp = regexp.MustCompile("^:([[:alpha:]][[:alnum:]]*):$")

func CreatePathTree() PathTree {
	return createPathNode("/")
}

func (rootNode *PathNode) AddPathHandler(
	method string,
	path string,
	handler Handler,
) PathTree {
	segments := extractPathSegments(path)
	addHandler(rootNode, segments, method, handler)
	return rootNode
}

func (rootNode *PathNode) GetPathContext(
	method string,
	path string,
) PathContext {
	segments := extractPathSegments(path)
	return createPathContext(rootNode, method, segments)
}

func createPathContext(
	node *PathNode,
	method string,
	segments []string,
) PathContext {
	context := PathContext{
		handlers:      make([]Handler, 0),
		pathVariables: make(map[string]string),
	}

	if isVariable(node.segment) {
		context.pathVariables[getVariableName(node.segment)] = segments[0]
	}

	remainingSegments := segments[1:]
	if len(remainingSegments) == 0 {
		context.handlers = append(context.handlers, node.pathHandlers[method]...)
	} else {
		children := make([]*PathNode, 0)
		child, exist := node.children[remainingSegments[0]]
		if exist {
			children = append(children, child)
		}
		variableChildren := getVariableChildren(node.children)
		children = append(children, variableChildren...)
		for _, child1 := range children {
			pathContext := createPathContext(child1, method, remainingSegments)
			context.handlers = append(context.handlers, pathContext.handlers...)
			for key, value := range pathContext.pathVariables {
				context.pathVariables[key] = value
			}
		}
	}

	return context
}

func addHandler(
	node *PathNode,
	segments []string,
	method string,
	handler Handler,
) {
	remainingSegments := append(segments[:0], segments[1:]...)
	if len(remainingSegments) == 0 {
		node.pathHandlers[method] = append(node.pathHandlers[method], handler)
	} else {
		segment := remainingSegments[0]
		childNode, exist := node.children[segment]
		if !exist {
			panicIfInvalid(node, segment)
			childNode = createPathNode(segment)
		}
		node.children[segment] = childNode
		addHandler(childNode, remainingSegments, method, handler)
	}
}
func isVariable(segment string) bool {
	return varRegexp.MatchString(segment)
}

func getVariableName(segment string) string {
	name := varRegexp.FindStringSubmatch(segment)
	return name[1]
}

func getVariableChildren(nodes map[string]*PathNode) []*PathNode {
	possibleChildren := make([]*PathNode, 0)
	for segment := range nodes {
		if isVariable(segment) {
			possibleChildren = append(possibleChildren, nodes[segment])
		}
	}

	return possibleChildren
}

func extractPathSegments(
	path string,
) []string {
	segments := strings.Split(path, "/")
	segments = filter(segments, func(element string) bool {
		return len(element) > 0
	})

	return append([]string{"/"}, segments...)
}

func filter(
	list []string,
	filterFunc func(string) bool,
) []string {
	filtered := make([]string, 0)
	for _, element := range list {
		if filterFunc(element) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}

func createPathNode(
	segment string,
) *PathNode {
	return &PathNode{
		segment:      segment,
		children:     make(map[string]*PathNode),
		pathHandlers: make(map[string][]Handler),
	}
}

func panicIfInvalid(node *PathNode, segment string) {
	for key := range node.children {
		if isVariable(key) || isVariable(segment) && key != segment {
			panic(fmt.Sprintf("Invalid config. Parent: %s, "+
				"conflicting path segments: %s, %s", node.segment, key, segment))
		}
	}
}
