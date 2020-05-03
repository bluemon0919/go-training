package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("http://localhost:8000/ ...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	v := r.URL.Query()
	if v == nil {
		return
	}
	for key, vs := range v {
		fmt.Printf("%s = %s\n", key, vs[0])
		fmt.Fprintf(w, "%s = %s\n", key, vs[0])
	}
}
