package main

import (
	"log"
	"net/http"

	"github.com/aargeee/whwh/handlers"
)

func main() {
	log.Fatal(
		http.ListenAndServe(":8000", handlers.NewServer()),
	)
}
