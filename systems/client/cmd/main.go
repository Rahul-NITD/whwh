package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Rahul-NITD/whwh/systems/client"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	serverUrl, ok := viper.Get("TESTERSERVER_URL").(string)
	if !ok {
		log.Fatalln("Could not get TESTERSERVER_URL")
	}

	hookUrl := flag.String("h", "http://localhost:3000", "specify the hook url, defaults to localhost:3000")

	c, sid, err := client.ClientConnect(serverUrl, *hookUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("All data to %s?stream=%s will be forwarded to %s", serverUrl, sid, *hookUrl)
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
