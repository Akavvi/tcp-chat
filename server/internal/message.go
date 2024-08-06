package internal

import (
	"fmt"
	"time"
)

type Message struct {
	Text   string
	Time   time.Time
	Client *Client
}

func NewMessage(text string, time time.Time, client *Client) *Message {
	return &Message{
		Text:   text,
		Time:   time,
		Client: client,
	}
}

func (m Message) String() string {
	return fmt.Sprintf("%s - %s: %s\n", m.Time.Format(time.Kitchen), m.Client.Name, m.Text)
}
