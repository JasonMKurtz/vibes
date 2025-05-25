package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	Logger      = log.Default()
	ErrorLogger = log.Default()
)

// Client represents a connected IRC client.
type Client struct {
	Conn     net.Conn
	Nickname string
	Username string
	Channels map[string]bool
}

// Server maintains IRC state.
type Server struct {
	Addr     string
	ln       net.Listener
	mu       sync.Mutex
	clients  map[net.Conn]*Client
	channels map[string]map[*Client]bool
}

// NewServer creates a new IRC server.
func NewServer(addr string) *Server {
	return &Server{
		Addr:     addr,
		clients:  make(map[net.Conn]*Client),
		channels: make(map[string]map[*Client]bool),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln
	s.Addr = ln.Addr().String()
	Logger.Printf("IRC server listening on %s", s.Addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				return nil
			}
			ErrorLogger.Println("accept error:", err)
			continue
		}
		Logger.Printf("Client connected: %s", conn.RemoteAddr())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	client := &Client{Conn: conn, Channels: make(map[string]bool)}
	s.mu.Lock()
	s.clients[conn] = client
	s.mu.Unlock()

	defer func() {
		Logger.Printf("Client disconnected: %s", conn.RemoteAddr())
		s.mu.Lock()
		for ch := range client.Channels {
			delete(s.channels[ch], client)
			if len(s.channels[ch]) == 0 {
				delete(s.channels, ch)
			}
			Logger.Printf("%s left %s", client.Nickname, ch)
		}
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
		c.Nickname = arg
	case "USER":
		c.Username = arg
	case "PING":
		c.Conn.Write([]byte("PONG :" + arg + "\r\n"))
	case "JOIN":
		s.joinChannel(c, arg)
	case "PRIVMSG":
		s.handlePrivMsg(c, arg)
	case "QUIT":
		c.Conn.Close()
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
	if c.Channels == nil {
		c.Channels = make(map[string]bool)
	}
	c.Channels[name] = true
	s.mu.Unlock()
	Logger.Printf("%s joined %s", c.Nickname, name)
	s.broadcast(ch, fmt.Sprintf(":%s JOIN %s\r\n", c.Nickname, name))
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
		s.broadcast(ch, fmt.Sprintf(":%s PRIVMSG %s :%s\r\n", c.Nickname, target, body))
	} else {
		s.mu.Lock()
		for client := range s.clients {
			if s.clients[client].Nickname == target {
				client.Write([]byte(fmt.Sprintf(":%s PRIVMSG %s :%s\r\n", c.Nickname, target, body)))
				break
			}
		}
		s.mu.Unlock()
	}
}

func (s *Server) broadcast(clients map[*Client]bool, msg string) {
	for c := range clients {
		c.Conn.Write([]byte(msg))
	}
}

// Close shuts down the server listener.
func (s *Server) Close() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}
