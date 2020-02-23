package internal

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn    net.Conn
	Id      int
	Friends []int
	Server  TCPServer
}

func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onConnectionClosed(c, err)
			return
		}
		fmt.Println(c.Server)
		c.Server.onNewMessage(c, message)
	}
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
