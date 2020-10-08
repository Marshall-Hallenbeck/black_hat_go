package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	r1, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatalln("Error with GET")
	}
	fmt.Println(r1.Status)

	body, err := ioutil.ReadAll(r1.Body)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(string(body))
	r1.Body.Close()

	r2, err := http.Head("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatalln("Error with HEAD")
	}
	// Read response body. Not shown.
	defer r2.Body.Close()
	form := url.Values{}
	form.Add("foo", "bar")
	r3, err := http.Post(
		"https://www.google.com/robots.txt",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		log.Fatalln("Error with POST")
	}
	// Read response body. Not shown.
	defer r3.Body.Close()
}