package main

import (
	"log"
	"net"
	"tcp-chat/server/internal"
	"tcp-chat/server/pkg"
)

func main() {
	hub := internal.NewHub()

	conn, err := net.Listen(pkg.TYPE, pkg.PORT)
	if err != nil {
		panic(err)
	}
	defer func(conn net.Listener) {
		_ = conn.Close()
	}(conn)
	log.Println("Listening on " + pkg.PORT)

	for {
		client, err := conn.Accept()
		if err != nil {
			panic(err)
		}
		hub.JoinClient(internal.NewClient(client))
	}
}
