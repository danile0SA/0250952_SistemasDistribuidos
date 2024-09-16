package main

import (
	"log"

	server "github.com/danile0_SA/0250952_SistemasDistribuidos/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
