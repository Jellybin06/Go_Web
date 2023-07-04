package main

import (
	"net/http"
	"workspace_01/hello_world/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
