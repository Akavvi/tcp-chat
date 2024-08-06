package internal

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	Name     string
	Room     *Room
	Incoming chan *Message
	Outgoing chan string
	Conn     net.Conn
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	LoggedIn bool
	Logining bool
}

func NewClient(conn net.Conn) *Client {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)
	client := &Client{
		Name:     "",
		Room:     nil,
		Incoming: make(chan *Message),
		Outgoing: make(chan string),
		Conn:     conn,
		Reader:   reader,
		Writer:   writer,
		LoggedIn: false,
		Logining: false,
	}

	client.Listen()

	return client
}

func (c Client) Close() {
	_ = c.Conn.Close()
}

func (c Client) Listen() {
	go c.Read()
	go c.Write()
}

func (c *Client) Write() {
	for m := range c.Outgoing {
		_, err := c.Writer.WriteString(m + "\n")
		if err != nil {
			log.Printf("Error writing to client: %v\n", err)
			break
		}
		err = c.Writer.Flush()
		if err != nil {
			log.Printf("Error flushing to client: %v\n", err)
			break
		}
	}
	log.Println("Client writer disconnected")
}

func (c *Client) Read() {
	for {
		text, err := c.Reader.ReadString('\n')
		if err != nil {
			break
		}
		message := NewMessage(strings.TrimSuffix(text, "\n"), time.Now(), c)

		c.Incoming <- message
	}
	close(c.Incoming)
	log.Println("Client reader disconnected")
}
