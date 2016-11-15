// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"crypto/tls"
	"log"
	"net"

	"strings"

	"github.com/DanielOaks/girc-go/ircmsg"
)

// Client is a client connection (shared between server tests)
type Client struct {
	Nick   string
	Socket *Socket
}

// Connect connects to the given server
func (c *Client) Connect(s *Server) error {
	// connect
	var conn net.Conn
	var err error

	if s.Conn.IsTLS {
		conn, err = tls.Dial("tcp", s.Conn.Address, nil)
	} else {
		conn, err = net.Dial("tcp", s.Conn.Address)
	}
	if err != nil {
		log.Fatal("Could not connect:", err.Error())
		return err
	}

	// create socket
	socket := NewSocket(conn)
	c.Socket = &socket

	return nil
}

// Disconnect disconnects from the given server
func (c *Client) Disconnect(s *Server) {
	if c.Socket.Closed {
		log.Println("Disconnected early")
		//TODO(dan): mark as closed badly
	} else {
		//DEBUG(dan): log.Println(c.Nick, "disconnecting")
		c.Socket.WriteLine("QUIT\r\n")
		// wait 'til we get ERROR message back
		for {
			line, err := c.Socket.Read()
			if err != nil {
				log.Println("Disconnected incorrectly 1:", err.Error())
				//TODO(dan): mark as closed badly
				break
			}
			//DEBUG(dan): log.Println(c.Nick, "got line:", strings.TrimSpace(line))

			msg, err := ircmsg.ParseLine(line)
			if err != nil {
				log.Println("Disconnected incorrectly 2:", err.Error())
				//TODO(dan): mark as closed badly
				break
			}

			if strings.ToUpper(msg.Command) == "ERROR" {
				//TODO(dan): mark as closed nicely
				break
			}
		}
	}
	c.Socket = nil
}
