package irc

import (
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(":6667")
	if s.Addr != ":6667" {
		t.Errorf("expected addr to be :6667, got %s", s.Addr)
	}
	if s.clients == nil || s.channels == nil {
		t.Error("expected maps to be initialized")
	}
}

func TestBroadcast(t *testing.T) {
	s := NewServer(":0")
	client := &Client{}
	conn1, conn2 := net.Pipe()
	client.Conn = conn1
	ch := map[*Client]bool{client: true}

	go s.broadcast(ch, "test\r\n")

	buf := make([]byte, 6)
	if _, err := conn2.Read(buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "test\r\n" {
		t.Errorf("expected message 'test\\r\\n', got %q", string(buf))
	}
}
