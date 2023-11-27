package main

import (
	"io"
	"net/http"
	"net/url"
)

var CallBackURL *url.URL

func handleCallback(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Close this page")
	CallBackURL = r.URL
}

func main() {
	LoginExample("http://localhost:8200")
}
