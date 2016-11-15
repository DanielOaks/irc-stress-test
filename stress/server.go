// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

// ServerConnectionDetails holds the details used to connect to the server.
type ServerConnectionDetails struct {
	Address string
	IsTLS   bool
}

// Server represents a server we are stress-testing.
type Server struct {
	Name string
	Conn ServerConnectionDetails
}
