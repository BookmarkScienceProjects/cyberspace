package main

import (
	"github.com/stojg/vivere/lib/client"
	"golang.org/x/net/websocket"
	"net/http"
	"sync"
)

type clientHub struct {
	sync.Mutex
	clients []*client.Client
}

func (network *clientHub) add(c *client.Client) {
	network.Lock()
	defer network.Unlock()
	network.clients = append(network.clients, c)
}

func (network *clientHub) Write(data []byte) (n int, err error) {
	network.Lock()
	defer network.Unlock()
	for _, client := range network.clients {
		nc, err := client.Update(data, 1)
		n += nc
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func initNetwork() *clientHub {
	Println("Inititalising Network")

	hub := &clientHub{
		clients: make([]*client.Client, 0),
	}
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

	go func(client chan *client.Client, h *clientHub) {
		for newClient := range client {
			Println("New client connected")
			h.add(newClient)
		}
	}(ch.NewClients(), hub)

	go func() {
		Println(http.ListenAndServe(":8080", nil))
	}()

	return hub

}
