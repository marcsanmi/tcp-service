package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	address            string
	clients            map[int]*Client
	onNewMessage       func(c *Client, message string)
	onConnectionClosed func(c *Client, err error)
}

type Message struct {
	Id      int   `json:"id"`
	Friends []int `json:"friends"`
}

func NewServer(addr string) *TCPServer {
	return &TCPServer{
		address: addr,
		clients: make(map[int]*Client),
	}
}

func (s *TCPServer) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	// Listen for connections
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}
		client := &Client{
			conn:   conn,
			Server: *s,
		}
		// Handle client request
		go client.listen()
	}
}

func (s *TCPServer) Close() error {
	return nil
}

func (s *TCPServer) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

func (s *TCPServer) InitListeners() {
	// Init all the listeners
	s.OnNewMessage(func(c *Client, message string) {
		fmt.Println("! NEWE MESSAGE")
		var msg Message
		err := json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("out")
			c.Close()
			return
		}
		if err := c.Server.AddClient(c, msg); err != nil {
			c.Send(err.Error())
			c.Close()
			return
		}
		c.Server.NotifyFriendsWithMessage(c, map[string]bool{"online": true})
	})

	s.OnConnectionClosed(func(c *Client, err error) {
		c.Server.DeleteClient(c)
		c.Server.NotifyFriendsWithMessage(c, map[string]bool{"online": false})
	})
}

func (s *TCPServer) OnConnectionClosed(callback func(c *Client, err error)) {
	s.onConnectionClosed = callback
}

func (s *TCPServer) AddClient(client *Client, message Message) error {
	client.Id = message.Id
	client.Friends = message.Friends
	_, ok := s.clients[client.Id]
	fmt.Println("adding client..")
	if ok {
		fmt.Println("BYEE!")
		return errors.New("User id already exists")
	}
	s.clients[client.Id] = client
	return nil
}

func (s *TCPServer) DeleteClient(c *Client) {
	delete(s.clients, c.Id)
}

func (s *TCPServer) NotifyFriendsWithMessage(c *Client, message map[string]bool) {
	response, err := json.Marshal(message)
	if err != nil {
		response = nil
	}
	fmt.Println(s.clients)
	for _, client := range s.clients {
		for _, id := range client.Friends {
			if id == c.Id {
				client.SendBytes(response)
			}
		}
	}
}
