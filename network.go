package main

import (
	"github.com/stojg/vivere/lib/client"
	"golang.org/x/net/websocket"
	"net/http"
)

var clients []*client.Client

func init() {
	Println("Inititalising Network")

	ch := client.NewClientManager()
	http.Handle("/ws/", websocket.Handler(ch.Websocket))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		http.ServeFile(w, r, "static/"+r.URL.Path[1:])
	})

	go func() {
		Println(http.ListenAndServe(":8080", nil))
	}()

	go func(client chan *client.Client) {
		for newClient := range client {
			Println("New client connected")
			clients = append(clients, newClient)
		}
	}(ch.NewClients())
}
