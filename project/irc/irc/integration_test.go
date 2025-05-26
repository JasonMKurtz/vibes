package irc

import (
	"errors"
	"net"
	"strings"
	"testing"

	ic "vibes/client"
)

func TestClientFlow(t *testing.T) {
	s := NewServer(":0")
	go func() {
		if err := s.Run(); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Fatalf("server error: %v", err)
		}
	}()
	<-s.Ready()
	defer s.Close()

	c1, err := ic.Connect(s.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Close()

	c2, err := ic.Connect(s.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Close()

	if err := c1.Login("alice"); err != nil {
		t.Fatal(err)
	}
	if err := c2.Login("bob"); err != nil {
		t.Fatal(err)
	}

	c1.Join("#room")
	c2.Join("#room")

       // read join messages: one for c1 and two after c2 joins
       c1.ReadLine()
       c2.ReadLine()
       c1.ReadLine()

	c1.Msg("#room", "hello")

	received := false
	for i := 0; i < 5; i++ {
		line, err := c2.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(line, "PRIVMSG #room :hello") {
			received = true
			break
		}
	}
	if !received {
		t.Fatal("privmsg not received")
	}

	c2.Part("#room")
	partReceived := false
	for i := 0; i < 5; i++ {
		line, err := c1.ReadLine()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(line, "PART #room") {
			partReceived = true
			break
		}
	}
	if !partReceived {
		t.Fatal("part not received")
	}
}
