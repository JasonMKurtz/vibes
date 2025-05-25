package main

import (
	"log"
	"os"
	"time"

	"vibes/client"
	"vibes/irc"
)

func main() {
	irc.Logger = log.New(os.Stdout, "", log.LstdFlags)
	irc.ErrorLogger = irc.Logger

	srv := irc.NewServer(":6667")
	ready := make(chan struct{})
	go func() {
		if err := srv.Run(ready); err != nil {
			irc.ErrorLogger.Fatal(err)
		}
	}()

	// Wait for the server to start
	<-ready

	cli, err := client.Connect("localhost:6667")
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	cli.Login("tester")
	cli.Join("#chat")
	cli.Msg("#chat", "hello world")

	time.Sleep(100 * time.Millisecond)
}
