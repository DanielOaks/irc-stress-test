// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import "sync/atomic"

// ServerConnectionDetails holds the details used to connect to the server.
type ServerConnectionDetails struct {
	Address string
	IsTLS   bool
}

// Server represents a server we are stress-testing.
type Server struct {
	// stats
	succeeded uint64 // align to 64-bit boundary

	Name string
	Conn ServerConnectionDetails
}

func (server *Server) RecordSuccess() {
	atomic.AddUint64(&server.succeeded, 1)
}

func (server *Server) Succeeded() uint64 {
	return atomic.LoadUint64(&server.succeeded)
}
