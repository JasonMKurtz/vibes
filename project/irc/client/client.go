package client

import (
	"bufio"
	"fmt"
	"net"
)

// Client provides helper methods for IRC interactions.
type Client struct {
	conn net.Conn
	r    *bufio.Reader
}

// Connect establishes a connection to the IRC server.
func Connect(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, r: bufio.NewReader(conn)}, nil
}

// Nick sends the NICK command.
func (c *Client) Nick(name string) error {
	return c.sendf("NICK %s", name)
}

// User sends the USER command.
func (c *Client) User(user string) error {
	return c.sendf("USER %s", user)
}

// Join joins the given channel.
func (c *Client) Join(channel string) error {
	return c.sendf("JOIN %s", channel)
}

// Msg sends a PRIVMSG to the target.
func (c *Client) Msg(target, message string) error {
	return c.sendf("PRIVMSG %s :%s", target, message)
}

// ReadLine reads a line from the server.
func (c *Client) ReadLine() (string, error) {
	line, err := c.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) sendf(format string, args ...any) error {
	_, err := fmt.Fprintf(c.conn, format+"\r\n", args...)
	return err
}
