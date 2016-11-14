// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

// Server represents a server we are stress-testing.
type Server struct {
	Name  string
	Addr  string
	IsTLS bool
}
