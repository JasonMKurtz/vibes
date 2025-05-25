package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"vibes/client"
)

func main() {
	addr := flag.String("server", "localhost:6667", "IRC server address")
	flag.Parse()

	c, err := client.Connect(*addr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Println("Connected to", *addr)

	joined := make(map[string]bool)
	nickname := ""

	go func() {
		for {
			line, err := c.ReadLine()
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			line = strings.TrimSpace(line)
			parts := strings.SplitN(line, " :", 2)
			msg := ""
			if len(parts) == 2 {
				msg = parts[1]
			}
			fields := strings.Fields(parts[0])
			if len(fields) < 2 {
				continue
			}
			src := strings.TrimPrefix(fields[0], ":")
			cmd := fields[1]
			switch cmd {
			case "JOIN":
				if len(fields) >= 3 {
					ch := fields[2]
					if src == nickname {
						joined[ch] = true
					}
					if joined[ch] {
						fmt.Printf("%s joined %s\n", src, ch)
					}
				}
			case "PART":
				if len(fields) >= 3 {
					ch := fields[2]
					if joined[ch] {
						fmt.Printf("%s left %s\n", src, ch)
						if src == nickname {
							delete(joined, ch)
						}
					}
				}
			case "PRIVMSG":
				if len(fields) >= 3 {
					target := fields[2]
					if joined[target] {
						fmt.Printf("[%s] %s: %s\n", target, src, msg)
					}
				}
			}
		}
	}()

	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			fmt.Print("> ")
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		cmd := strings.ToLower(parts[0])
		arg := ""
		if len(parts) > 1 {
			arg = parts[1]
		}
		switch cmd {
		case "login":
			if arg == "" {
				fmt.Println("usage: login <name>")
			} else if err := c.Login(arg); err != nil {
				fmt.Println("error:", err)
			} else {
				nickname = arg
			}
		case "join":
			if arg == "" {
				fmt.Println("usage: join <channel>")
			} else if err := c.Join(arg); err != nil {
				fmt.Println("error:", err)
			}
		case "part":
			if arg == "" {
				fmt.Println("usage: part <channel>")
			} else if err := c.Part(arg); err != nil {
				fmt.Println("error:", err)
			}
		case "msg":
			args := strings.SplitN(arg, " ", 2)
			if len(args) != 2 {
				fmt.Println("usage: msg <target> <message>")
			} else if err := c.Msg(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
			}
		case "read":
			line, err := c.ReadLine()
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Print(line)
			}
		case "quit":
			c.Close()
			return
		default:
			fmt.Println("unknown command:", cmd)
		}
		fmt.Print("> ")
	}
}
