package main

import (
	"log"
	"net/http"

	"github.com/Rahul-NITD/whwh/systems/hook"
)

func main() {
	hook := hook.NewHook(&hook.OsStdoutAdapter{})
	log.Fatal(http.ListenAndServe("localhost:3000", hook))
}
