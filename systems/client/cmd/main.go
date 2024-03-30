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
	serverUrl = "http://localhost:8000"
	if os.Getenv("AARGEEE_PROD_ENV") == "True" {
		serverUrl = os.Getenv("AARGEEE_WHWH_SERVER")
	} else {
		println("Running in test environment, ", serverUrl)
	}
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
