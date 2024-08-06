package internal

import (
	"log"
	"strings"
	"tcp-chat/server/pkg"
)

type Hub struct {
	Clients  []*Client
	Rooms    map[string]*Room
	Incoming chan *Message
	Join     chan *Client
	Leave    chan *Client
}

func NewHub() *Hub {
	hub := &Hub{
		Clients:  make([]*Client, 0),
		Rooms:    make(map[string]*Room),
		Incoming: make(chan *Message),
		Join:     make(chan *Client),
		Leave:    make(chan *Client),
	}
	hub.Listen()

	return hub
}

func (h Hub) Listen() {
	go func() {
		for {
			select {
			case msg := <-h.Incoming:
				h.HandleMessage(msg)
			case client := <-h.Join:
				h.JoinClient(client)
			case client := <-h.Leave:
				h.LeaveClient(client)
			}
		}
	}()
}
func (h Hub) HandleMessage(message *Message) {
	switch {
	default:
		SendMessage(message)
	case !message.Client.LoggedIn:
		Login(message.Client, strings.TrimSpace(message.Text))
	case strings.HasPrefix(message.Text, pkg.CREATE):
		name := strings.TrimSpace(strings.TrimPrefix(message.Text, pkg.CREATE+" "))
		Create(&h, message.Client, name)
	case strings.HasPrefix(message.Text, pkg.JOIN):
		name := strings.TrimSpace(strings.TrimPrefix(message.Text, pkg.JOIN+" "))
		Join(&h, message.Client, name)
	case strings.HasPrefix(message.Text, pkg.LEAVE):
		Leave(message.Client)
	case strings.HasPrefix(message.Text, pkg.HELP):
		Help(message.Client)
	case strings.HasPrefix(message.Text, pkg.LIST):
		List(&h, message.Client)
	}
}

func (h *Hub) JoinClient(client *Client) {
	h.Clients = append(h.Clients, client)
	client.Outgoing <- "Welcome to the chat!"
	go func() {
		for msg := range client.Incoming {
			h.Incoming <- msg
		}
		h.Leave <- client
	}()
}

func (h *Hub) LeaveClient(client *Client) {
	if client.Room != nil {
		client.Room.Leave(client)
	}
	for i, clients := range h.Clients {
		if clients == client {
			h.Clients = append(h.Clients[:i], h.Clients[i+1:]...)
		}
	}
	close(client.Outgoing)
	log.Println("Client outgoing closed")
}

func (h *Hub) SendMessage(message *Message) {
	for _, client := range h.Clients {
		if client != message.Client {
			client.Outgoing <- message.Text
		}
	}
}
