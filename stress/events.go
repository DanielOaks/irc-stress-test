// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
)

// EventQueue represents a series of events.
type EventQueue struct {
	Events   []Event
	Finished chan bool
}

// NewEventQueue returns a new EventQueue
func NewEventQueue() EventQueue {
	events := EventQueue{
		Events:   make([]Event, 0),
		Finished: make(chan bool),
	}
	return events
}

// Run goes through our event list.
func (queue EventQueue) Run(server *Server, clients map[int]*Client) {
	for _, event := range queue.Events {
		client := clients[event.Client]

		if event.Type == ETConnect {
			err := client.Connect(server)
			if err != nil {
				log.Fatal("Could not connect...", err.Error())
			}
		} else if event.Type == ETDisconnect {
			client.Disconnect(server)
			client.Socket = nil
		} else if event.Type == ETLine {
			line := *event.Line
			client.Socket.Write(line)
		} else if event.Type == ETWait {
			log.Println("ETWait events not yet implemented")
		} else {
			log.Println("Got unknown event type:", event.Type)
			spew.Dump(event)
			fmt.Print("\n\n")
		}
	}
	// send finished notice, used for syncing
	queue.Finished <- true
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
	Client int
	Type   EventType
	Line   *string
	Wait   *WaitMessage
}
