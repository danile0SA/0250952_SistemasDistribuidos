package main

import (
	"net/http"
)

func main() {

	srv := server_api.NewServer()
	http.ListenAndServe(":8080", srv)
}
