package main

import (
	"github.com/gorilla/websocket"

	r "github.com/dancannon/gorethink"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type FindHandler func(string) (Handler, bool)

type Client struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
	session     *r.Session
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.socket.Close()
}

func (client *Client) Write() {
	for msg := range client.send {
		if err := client.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.socket.Close()
}

func NewClient(scoket *websocket.Conn, findHandler FindHandler, session *r.Session) *Client {
	return &Client{
		send:        make(chan Message),
		socket:      scoket,
		findHandler: findHandler,
		session:     session,
	}
}
