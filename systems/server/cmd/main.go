package main

import (
	"log"
	"net/http"

	"github.com/Rahul-NITD/whwh/handlers"
)

func main() {
	log.Fatal(
		http.ListenAndServe(":8000", handlers.NewTesterServerHandler()),
	)
}
