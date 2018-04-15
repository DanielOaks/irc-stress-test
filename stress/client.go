// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/goshuirc/irc-go/ircmsg"

	"strings"
)

// Client is a client connection
type Client struct {
	Nick   string
	Socket *Socket
}

func NewClient(id int) Client {
	return Client{
		Nick: fmt.Sprintf("ircstress_%d", id),
	}
}

// Connect connects to the given server
func (c *Client) Connect(server *Server) error {
	// connect
	var conn net.Conn
	var err error

	addr := strings.TrimPrefix(server.Conn.Address, "unix:")

	if server.Conn.IsTLS {
		conn, err = tls.Dial("tcp", addr, nil)
	} else if strings.HasPrefix(addr, "/") {
		conn, err = net.Dial("unix", addr)
	} else {
		conn, err = net.Dial("tcp", server.Conn.Address)
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
func (c *Client) Disconnect(server *Server) {
	// issue #4: report to other clients that we are ready to disconnect
	server.ClientsReadyToDisconnect.Done()
	if c.Socket.Closed {
		log.Println("Disconnected early")
		//TODO(dan): mark as closed badly
	} else {
		// wait for everyone to else to report the same
		server.ClientsReadyToDisconnect.Wait()
		//DEBUG(dan): log.Println(c.Nick, "disconnecting")
		c.Socket.WriteLine("QUIT")
		// wait 'til we get ERROR message back
		var lastLine string
		var totalLines int
		for {
			line, err := c.Socket.Read()
			if err != nil {
				log.Println("Disconnected incorrectly 1:", err.Error())
				log.Println("last line:", totalLines, ":", lastLine)
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

			// fmt.Println(c.Nick, line)
			lastLine = line
			totalLines++

			if strings.ToUpper(msg.Command) == "ERROR" {
				//TODO(dan): mark as closed nicely
				server.RecordSuccess()
				break
			}
		}
	}
	c.Socket = nil
}
