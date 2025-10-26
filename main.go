package main

import (
	"log"
	"net/http"
	"os"
)

var appVersion = "1.0.0"

func main() {
	if err := run(os.Args[1:], os.Exit, http.ListenAndServe); err != nil {
		log.Fatal(err)
	}
}
