package main

import (
	"log"

	"github.com/tkircsi/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":5000")
	log.Fatal(srv.ListenAndServe())
}
