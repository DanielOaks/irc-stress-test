// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

// EventQueue represents a series of events.
type EventQueue []Event

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
	Client Client
	Type   EventType
	Line   *string
	Wait   *WaitMessage
}
