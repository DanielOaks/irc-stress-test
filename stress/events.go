// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"fmt"
	"log"
)

// EventQueue represents a series of events.
type EventQueue struct {
	Client Client
	Events []Event
	id     int
}

// NewEventQueue returns a new EventQueue
func NewEventQueue(id int) EventQueue {
	events := EventQueue{
		Events: make([]Event, 0),
		Client: NewClient(id),
		id:     id,
	}
	return events
}

// Run goes through our event list.
func (queue EventQueue) Run(server *Server) {
	client := &queue.Client
	for _, event := range queue.Events {
		switch event.Type {
		case ETConnect:
			fmt.Println("c", client.Nick)
			err := client.Connect(server)
			if err != nil {
				log.Fatal("Could not connect...", err.Error())
			}
		case ETDisconnect:
			client.Disconnect(server)
		case ETLine:
			client.Socket.Write(event.Line)
		case ETWait:
			log.Println("ETWait events not yet implemented")
		case ETPing:
			client.Ping()
		default:
			panic(fmt.Sprintf("Unknown event type: %d", event.Type))
		}
	}
	// send finished notice, used for syncing
	server.ClientsFinished.Done()
}

// EventType is the type of event it is.
type EventType int

const (
	// ETConnect represents the client connecting to the server.
	ETConnect EventType = iota
	// ETDisconnect represents the client disconnecting from the server.
	ETDisconnect
	// ETLine represents the client sending an IRC line to the server.
	ETLine
	// ETWait represents the client waiting for a specific response from the server.
	ETWait
	// ETPing causes the client to send a ping, then wait for the specific response
	ETPing
)

// WaitMessage is a message that the client should wait for.
type WaitMessage struct {
	// Command is the IRC command to wait for.
	Command *string
	// Params are the IRC message params to wait for.
	Params *string
}

// Event is an IRC event.
type Event struct {
	Type EventType
	Line string
	Wait *WaitMessage
}
