package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// Client represents a connected IRC client.
type Client struct {
	conn     net.Conn
	nickname string
	username string
	channels map[string]bool
}

// Server maintains IRC state.
type Server struct {
	addr     string
	mu       sync.Mutex
	clients  map[net.Conn]*Client
	channels map[string]map[*Client]bool
}

// NewServer creates a new IRC server.
func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		clients:  make(map[net.Conn]*Client),
		channels: make(map[string]map[*Client]bool),
	}
}

func (s *Server) run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Printf("IRC server listening on %s", s.addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	client := &Client{conn: conn, channels: make(map[string]bool)}
	s.mu.Lock()
	s.clients[conn] = client
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		line := reader.Text()
		s.handleLine(client, line)
	}
}

func (s *Server) handleLine(c *Client, line string) {
	parts := strings.SplitN(line, " ", 2)
	cmd := strings.ToUpper(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = parts[1]
	}
	switch cmd {
	case "NICK":
		c.nickname = arg
	case "USER":
		c.username = arg
	case "PING":
		c.conn.Write([]byte("PONG :" + arg + "\r\n"))
	case "JOIN":
		s.joinChannel(c, arg)
	case "PRIVMSG":
		s.handlePrivMsg(c, arg)
	case "QUIT":
		c.conn.Close()
	}
}

func (s *Server) joinChannel(c *Client, name string) {
	s.mu.Lock()
	ch, ok := s.channels[name]
	if !ok {
		ch = make(map[*Client]bool)
		s.channels[name] = ch
	}
	ch[c] = true
	c.channels[name] = true
	s.mu.Unlock()
	s.broadcast(ch, fmt.Sprintf(":%s JOIN %s\r\n", c.nickname, name))
}

func (s *Server) handlePrivMsg(c *Client, msg string) {
	parts := strings.SplitN(msg, " ", 2)
	if len(parts) != 2 {
		return
	}
	target, body := parts[0], parts[1]
	if strings.HasPrefix(body, ":") {
		body = body[1:]
	}
	if strings.HasPrefix(target, "#") {
		s.mu.Lock()
		ch := s.channels[target]
		s.mu.Unlock()
		s.broadcast(ch, fmt.Sprintf(
			":%s PRIVMSG %s :%s\r\n",
			c.nickname,
			target,
			body,
		))
	} else {
		s.mu.Lock()
		for client := range s.clients {
			if s.clients[client].nickname == target {
				client.Write([]byte(fmt.Sprintf(
					":%s PRIVMSG %s :%s\r\n",
					c.nickname,
					target,
					body,
				)))
				break
			}
		}
		s.mu.Unlock()
	}
}

func (s *Server) broadcast(clients map[*Client]bool, msg string) {
	for c := range clients {
		c.conn.Write([]byte(msg))
	}
}

func main() {
	addr := ":6667"
	s := NewServer(addr)
	if err := s.run(); err != nil {
		log.Fatal(err)
	}
}
