// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import(
	"github.com/rin1208/go-trace"
)
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	tracer trace.Tracer
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.tracer.Trace("新しいクライアントが参加したよ！")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.tracer.Trace("クライアントが退室したよ")
			}
		case message := <-h.broadcast:
			h.tracer.Trace("メッセージを受信したよ: ", string(message))
			for client := range h.clients {
				select {
				case client.send <- message:
				h.tracer.Trace(" -- クライアントに送信されたよ ")
				default:
					close(client.send)
					delete(h.clients, client)
					h.tracer.Trace("送信に失敗したよ")
				}
			}
		}
	}
}
