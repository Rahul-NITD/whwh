package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aargeee/whwh/systems/client"
)

func main() {
	serverUrl := "https://webhookwormhole-latest.onrender.com"
	hookUrl := flag.String("h", "http://localhost:3000", "specify the hook url, defaults to localhost:3000")
	flag.Parse()

	c, sid, err := client.ClientConnect(serverUrl, *hookUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("All data to %s?stream=%s will be forwarded to %q", serverUrl, sid, *hookUrl)
	unsubscribe, err := client.ClientSubscribe(c, sid, *hookUrl)
	if err != nil {
		log.Fatal(err)
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	unsubscribe()
	log.Println("Client Stopped")
}
