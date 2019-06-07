package main

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	r "github.com/dancannon/gorethink"
)

func addChannel(client *Client, data interface{}) {
	var channel Channel
	err := mapstructure.Decode(data, &channel)

	if err != nil {
		client.send <- Message{"error", err}
		return
	}

	go func() {
		r.Table("channel").Insert(channel).Exec(client.session)

		if err != nil {
			client.send <- Message{"error", err}
		}
	}()
}

func subscribeChannel(client *Client, data interface{}) {
	go func() {
		cursor, err := r.Table("channel").Changes(r.ChangesOpts{IncludeInitial: true}).Run(client.session)

		if err != nil {
			client.send <- Message{"error", err}
			return
		}
		var change r.ChangeResponse

		for cursor.Next(&change) {
			fmt.Println(change.NewValue)
			if change.NewValue != nil {
				client.send <- Message{"channel add", change.NewValue}
				fmt.Println("sent channel add message")
			}
		}

	}()
}
