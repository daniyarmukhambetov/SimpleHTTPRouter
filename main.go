package main

import (
	"fmt"
	"net/http"
	"time"
)

//examples
func f(rw http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Println(r.URL.Path)
	fmt.Println("f1", params) // path params
}

func f1(w http.ResponseWriter, r *http.Request, params map[string]string) {
	time.Sleep(10 * time.Second)
	fmt.Println("f2", params)
}
func main() {
	r := NewRouter("api")
	r.Get("v1/salam/:name", f)
	r.Get("v2/salam/:name", f)
	r.Get("v1/privet", f1)
	r.Get("v1/privet/:name", f1)
	http.ListenAndServe("localhost:8080", r)
}
