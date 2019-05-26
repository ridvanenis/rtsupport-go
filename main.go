package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// msgType, msg, err := socket.ReadMessage()
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		var inMessage Message
		var outMessage Message
		if err := socket.ReadJSON(&inMessage); err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("%#v\n", inMessage)

		switch inMessage.Name {
		case "channel add":
			_, err := addChannel(inMessage.Data)
			if err != nil {
				outMessage = Message{"err", err}
				if err := socket.WriteJSON(outMessage); err != nil {
					fmt.Println(err)
					break
				}
			}

		case "channel subscribe":
			go subscribeChannel(socket)
		}

		// fmt.Println(string(msg))

		// if err = socket.WriteMessage(msgType, msg); err != nil {
		// 	fmt.Println(err)
		// }

	}
}

func addChannel(data interface{}) (Channel, error) {
	var channel Channel

	err := mapstructure.Decode(data, &channel)

	if err != nil {
		return channel, err
	}

	channel.Id = "1"
	fmt.Println("add channel")
	return channel, nil
}

func subscribeChannel(socket *websocket.Conn) {
	//TODO: rethink Query / changefeed

	for {
		time.Sleep(time.Second * 1)
		message := Message{"channel add", Channel{"1", "Software Support"}}
		socket.WriteJSON(message)
		fmt.Println("sent new channel")
	}
}
