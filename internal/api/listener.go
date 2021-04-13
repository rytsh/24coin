package api

import (
	"github.com/google/uuid"
)

// record all clients here
var clients = map[uuid.UUID]chan []byte{}

func clientAdd(channel chan []byte) {
	clients[uuid.New()] = channel
}

func clientDelete(s uuid.UUID) {
	delete(clients, s)
}

// if cannot consume data, it is already closed
func BroadcastMessage(message []byte) {
	for client := range clients {
		select {
		case clients[client] <- message:
		default:
			clientDelete(client)
		}
	}
}
