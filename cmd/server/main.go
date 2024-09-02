package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	server "github.com/danile0SA/0250952_SistemasDistribuidos/Project_Module/Internal/Server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())

	// Ejemplo de inicialización del store
	f, err := os.OpenFile("logger.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	s := &store{
		File: f,
		buf:  bufio.NewWriter(f),
	}
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}

	// Ejemplo de Append
	data := []byte("Hola, mundo!")
	_, pos, err := s.Append(data)
	if err != nil {
		log.Fatal(err)
	}

	// Ejemplo de Read
	readData, err := s.Read(pos)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(readData))

	// Cerrar el store
	if err := s.Close(); err != nil {
		log.Fatal(err)
	}
}
