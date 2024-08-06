package internal

import (
	"fmt"
)

func SendMessage(message *Message) {
	if !message.Client.LoggedIn {
		return
	}
	if message.Client.Room == nil {
		message.Client.Outgoing <- "You are not in a room"
		return
	}
	message.Client.Room.HandleMessage(message.String())
}
func Help(client *Client) {
	client.Outgoing <- "Available commands: create, join, leave, help, name"
}

func List(hub *Hub, client *Client) {
	client.Outgoing <- "Available rooms: "
	for roomName := range hub.Rooms {
		client.Outgoing <- roomName
	}
}

func Create(hub *Hub, client *Client, roomName string) {
	if hub.Rooms[roomName] != nil {
		client.Outgoing <- "Room already exists"
		return
	}
	room := NewRoom(roomName)
	hub.Rooms[roomName] = room
	client.Outgoing <- fmt.Sprintf("Room %s was created", room.Name)
}

func Join(hub *Hub, client *Client, roomName string) {
	if hub.Rooms[roomName] == nil {
		client.Outgoing <- "Room does not exist"
		return
	}
	if client.Room != nil {
		client.Outgoing <- "You are already in a room"
		return
	}
	room := hub.Rooms[roomName]
	room.Join(client)
	client.Outgoing <- fmt.Sprintf("Joined room %s", roomName)
	room.HandleMessage(client.Name + " has joined the room")
}

func Leave(client *Client) {
	if client.Room == nil {
		client.Outgoing <- "You are not in a room"
		return
	}
	client.Room.Leave(client)
}

var loginingUsers = make(map[*Client]bool)

func Login(client *Client, username string) {
	if client.LoggedIn {
		client.Outgoing <- "You are already logged in"
		return
	}
	if !loginingUsers[client] {
		client.Outgoing <- "Enter your username"
		loginingUsers[client] = true
		return
	}
	client.Name = username
	client.LoggedIn = true
	client.Outgoing <- "You have logged in"
}
