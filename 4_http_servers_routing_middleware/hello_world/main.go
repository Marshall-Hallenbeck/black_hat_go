package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("name"))
}

func main() {
	http.HandleFunc("/hello", hello)
	// the second argument being nil means the package will use the DefaultServerMux as the underlying handler
	// there's also http.ListenAndServeTLS(), but that requires additional parameters
	http.ListenAndServe(":8000", nil)
}
