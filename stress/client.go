// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

// Client is a client connection (shared between server tests)
type Client struct {
	Nick   string
	Socket *Socket
}
