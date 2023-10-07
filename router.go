package goo

import (
	"log"
	"strings"
)

type router struct {
	roots map[string]*node
}

func NewRouter() *router {
	return &router{roots: make(map[string]*node)}
}

func parsePath(path string) []string {

	pp := strings.Split(path, "/")
	parts := make([]string, 0)

	for _, value := range pp {
		if value != "" {
			parts = append(parts, value)
			if value[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRouter(method string, path string, handlers []HandlerFunc) {

	log.Printf("Route %4s - %s", method, path)

	parts := parsePath(path)
	_, ok := r.roots[method]

	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(path, parts, 0, handlers)
}

func (r *router) getRouter(method string, path string) (*node, map[string]string) {

	params := make(map[string]string)
	searchParts := parsePath(path)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)

	if node != nil {
		parts := parsePath(node.path)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return node, params
	}

	return nil, nil

}

func (r *router) handler(context *Context) {
	node, params := r.getRouter(context.Methon, context.Path)
	if node != nil {
		context.Params = params
		context.handlers = node.handlers
		context.Next()
	} else {
		context.String(404, "404 NOT PATH: %s\n", context.Path)
	}
}
