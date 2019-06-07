package main

import (
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
)

type Channel struct {
	Id   string `json:"id" gorethink:"id,omitempty"`
	Name string `json:"name" gorethink:"name"`
}

type User struct {
	Id   string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}

func main() {

	session, err := r.Connect(r.ConnectOpts{
		Address:  "172.17.0.2:28015",
		Database: "rtsupport",
	})

	if err != nil {
		log.Panic(err)
		return
	}

	router := NewRouter(session)
	router.Handle("channel add", addChannel)
	router.Handle("channel subscribe", subscribeChannel)
	http.Handle("/", router)
	http.ListenAndServe(":4000", nil)
}
