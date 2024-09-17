package main_http

import (
	"log"

	server "github.com/danile0SA/0250952_SistemasDistribuidos/Internal/Server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
