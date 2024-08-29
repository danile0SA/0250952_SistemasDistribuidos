package main

import (
	"log"

	server "github.com/danile0SA/0250952_SistemasDistribuidos/Project_Module/Internal/Server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
