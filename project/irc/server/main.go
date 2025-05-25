package main

import (
	"io"
	"log"
	"os"

	"vibes/irc"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	irc.Logger = log.New(logFile, "", log.LstdFlags)
	irc.ErrorLogger = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR: ", log.LstdFlags)

	addr := ":6667"
	s := irc.NewServer(addr)
	if err := s.Run(nil); err != nil {
		irc.ErrorLogger.Fatal(err)
	}
}
