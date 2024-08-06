package internal

type Room struct {
	Name     string
	Clients  []Client
	Messages []string
}

func NewRoom(name string) *Room {
	return &Room{
		Name:     name,
		Clients:  make([]Client, 0),
		Messages: make([]string, 0),
	}
}

func (r Room) Join(client *Client) {
	client.Room = &r
	for _, msg := range r.Messages {
		client.Outgoing <- msg
	}
	r.Clients = append(r.Clients, *client)

}

func (r Room) Leave(client *Client) {
	r.HandleMessage(client.Name + " has left the room")
	for i, clients := range r.Clients {
		if clients == *client {
			r.Clients = append(r.Clients[:i], r.Clients[i+1:]...)
			break
		}
	}
	client.Room = nil
}

func (r Room) HandleMessage(message string) {
	r.Messages = append(r.Messages, message)

	for _, client := range r.Clients {
		client.Outgoing <- message
	}
}
