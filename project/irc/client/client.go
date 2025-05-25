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

// Login sends the NICK and USER commands using the same value.
// IRC servers typically require both commands during connection
// setup but in this client the values are always identical.  Login
// combines the two so callers don't have to issue them separately.
func (c *Client) Login(name string) error {
	if err := c.sendf("NICK %s", name); err != nil {
		return err
	}
	return c.sendf("USER %s", name)
}

// Join joins the given channel.
func (c *Client) Join(channel string) error {
	return c.sendf("JOIN %s", channel)
}

// Part parts the given channel.
func (c *Client) Part(channel string) error {
	return c.sendf("PART %s", channel)
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
