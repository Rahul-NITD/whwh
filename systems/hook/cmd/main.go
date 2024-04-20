package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aargeee/whwh/systems/hook"
)

func main() {

	hookPort := flag.Int("p", 3001, "specify the hook url, defaults to localhost:3000")
	flag.Parse()
	hook := hook.NewHook(&hook.OsStdoutAdapter{})
	println("Hook running on ", fmt.Sprintf("localhost:%d", *hookPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *hookPort), hook))
}
