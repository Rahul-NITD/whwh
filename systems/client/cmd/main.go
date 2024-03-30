package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aargeee/whwh/systems/client"
	"github.com/aargeee/whwh/systems/client/cli"
)

func main() {
	var serverUrl string
	if os.Getenv("AARGEEE_TEST_ENV") == "True" {
		serverUrl = "http://localhost:8000"
	} else {
		serverUrl = os.Getenv("AARGEEE_WHWH_SERVER_URL") // "https://webhookwormhole-latest.onrender.com"
	}
	println("Running in environment, ", serverUrl)

	hookUrl := cli.ParseFlags()

	c, sid, err := client.ClientConnect(serverUrl, hookUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("All data to %s?stream=%s will be forwarded to %q", serverUrl, sid, hookUrl)
	unsubscribe, err := client.ClientSubscribe(c, sid, hookUrl)
	if err != nil {
		log.Fatal(err)
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	unsubscribe()
	log.Println("Client Stopped")
}
