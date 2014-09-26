package main

import (
	"expvar"
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"runtime"
)

func numGoroutine() interface{} {
	return interface{}(runtime.NumGoroutine())
}

func main() {
	fmt.Println("Hello Go Sandbox!")

	r := mux.NewRouter()
	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})
	r.HandleFunc("/hoge/{id}", func (w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		fmt.Fprint(w, "hoge Hello id = ", vars["id"])
	})

	expvar.Publish("goroutineNum", expvar.Func(numGoroutine))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

