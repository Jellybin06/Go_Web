package main

import (
	"net/http"
	"github.com/go_web_projectR/RESTfulAPI/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}
