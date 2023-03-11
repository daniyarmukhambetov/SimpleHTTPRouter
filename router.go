package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getType(method string) int {
	switch method {
	case "GET":
		return 0
	case "POST":
		return 1
	case "PUT":
		return 2
	case "DELETE":
		return 3
	case "PATCH":
		return 4
	default:
		return -1
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type node struct {
	name     string
	childes  map[string]*node
	handlers []HandlerFunc
}

type handler struct {
	method      string
	handlerFunc HandlerFunc
}

func New(pref string) *node {
	return &node{name: pref, childes: make(map[string]*node), handlers: make([]HandlerFunc, 5)}
}

func (n *node) AddEndpoint(path string, handlerFunc handler) string {
	s := strings.SplitN(path, "/", 2)
	if len(s) == 1 {
		//for _, h := range handlerFuncs {
		if t := getType(handlerFunc.method); t == -1 {
			return "invalid method passed"
		}
		//}
		var handlers []HandlerFunc = make([]HandlerFunc, 5)
		handlers[getType(handlerFunc.method)] = handlerFunc.handlerFunc
		n.childes[s[0]] = &node{
			name:     s[0],
			handlers: handlers,
			childes:  make(map[string]*node),
		}
		return ""
	}
	_, ok := n.childes[s[0]]
	if ok == false {
		n.childes[s[0]] = &node{name: s[0], childes: make(map[string]*node), handlers: make([]HandlerFunc, 5)}
	}
	child := n.childes[s[0]]
	return child.AddEndpoint(s[1], handlerFunc)
}

func (n *node) FindPath(path, method string, params map[string]string) (HandlerFunc, string) {
	s := strings.SplitN(path, "/", 2)
	if n.name[0:1] == ":" {
		var res HandlerFunc = nil
		var err string = "not found"
		params[n.name[1:]] = s[0]
		if len(s) == 1 {
			if t := getType(method); n.handlers[t] == nil {
				return nil, "method unsupported"
			} else {
				return n.handlers[t], ""
			}
		}
		for _, child := range n.childes {
			res, err = child.FindPath(s[1], method, params)
			if res != nil {
				return res, err
			}
		}
		return res, err
	}
	if len(s) == 1 {
		if s[0] != n.name {
			return nil, "not found"
		} else {
			if t := getType(method); n.handlers[t] == nil {
				return nil, "method unsupported"
			} else {
				return n.handlers[t], ""
			}
		}
	}
	if n.name == s[0] {
		var res HandlerFunc = nil
		var err string = "not found"
		for _, child := range n.childes {
			res, err = child.FindPath(s[1], method, params)
			if res != nil {
				return res, err
			} else if err == "method unsupported" {
				return res, err
			}
		}
		return res, err
	} else {
		return nil, "not found"
	}
}

func (n *node) Represent() {
	fmt.Println(n.name, n.handlers)
	for _, ch := range n.childes {
		ch.Represent()
	}
}

type Router struct {
	root *node
}

func (r *Router) AddPath(path string, funcs handler) string {
	return r.root.AddEndpoint(path, funcs)
}

func (r *Router) FindPath(path string, method string) (HandlerFunc, string, map[string]string) {
	params := make(map[string]string)
	res, err := r.root.FindPath(path, method, params)
	return res, err, params
}

func (r *Router) Get(path string, f HandlerFunc) string {
	return r.AddPath(path, handler{method: "GET", handlerFunc: f})
}

func (r *Router) Post(path string, f HandlerFunc) string {
	return r.AddPath(path, handler{method: "POST", handlerFunc: f})
}
func (r *Router) Put(path string, f HandlerFunc) string {
	return r.AddPath(path, handler{method: "PUT", handlerFunc: f})
}

func (r *Router) Patch(path string, f HandlerFunc) string {
	return r.AddPath(path, handler{method: "PATCH", handlerFunc: f})
}

func (r *Router) Delete(path string, f HandlerFunc) string {
	return r.AddPath(path, handler{method: "DELETE", handlerFunc: f})
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h, err, params := r.FindPath(req.URL.Path[1:], req.Method)
	fmt.Println(req.URL.Path)
	if h == nil {
		fmt.Println(err)
		return
	}
	h(rw, req, params)
}

func NewRouter(name string) *Router {
	return &Router{
		root: New(name),
	}
}
