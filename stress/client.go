// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"

	//"github.com/goshuirc/irc-go/ircmsg"

	"strconv"
	"strings"
)

var skipVerifyConfig = &tls.Config{
	InsecureSkipVerify: true,
}

// Client is a client connection
type Client struct {
	sync.Mutex

	Nick   string
	Socket *Socket
	closed chan bool

	pongEvent chan bool

	closeExpected bool
	pingCounter   uint64
	lastPong      uint64
	lastLine      string
	totalLines    int
}

func NewClient(id int) Client {
	return Client{
		Nick:        fmt.Sprintf("ircstress_%d", id),
		closed:      make(chan bool, 1),
		pongEvent:   make(chan bool, 1),
		pingCounter: 1,
	}
}

func (client *Client) SetCloseExpected(val bool) {
	client.Lock()
	defer client.Unlock()
	client.closeExpected = val
}

func (client *Client) CloseExpected() bool {
	client.Lock()
	defer client.Unlock()
	return client.closeExpected
}

func (client *Client) recordPong(pong uint64) {
	client.Lock()
	defer client.Unlock()
	if pong > client.lastPong {
		client.lastPong = pong
	}
}

func (client *Client) LastPong() uint64 {
	client.Lock()
	defer client.Unlock()
	return client.lastPong
}

func (client *Client) Ping() {
	client.Lock()
	ping := client.pingCounter
	client.pingCounter++
	client.Unlock()

	client.Socket.Write(fmt.Sprintf("PING %d\r\n", ping))
	for {
		<-client.pongEvent
		if client.LastPong() >= ping {
			return
		}
	}
}

func (client *Client) readLoop(server *Server) {
	quitRecvd := false

	for {
		line, err := client.Socket.Read()
		if err != nil && !quitRecvd {
			log.Println("Disconnected incorrectly 1:", err.Error())
			log.Println("last line:", client.totalLines, ":", client.lastLine)
			//TODO(dan): mark as closed badly
		}
		if err != nil {
			break
		}

		if strings.HasPrefix(line, "ERROR Quit") {
			if client.CloseExpected() {
				server.RecordSuccess()
				quitRecvd = true
			} else {
				log.Println(client.Nick, "unexpected quit")
			}
		} else {
			pieces := strings.Split(line, " ")
			if len(pieces) > 1 && pieces[len(pieces)-2] == "PONG" {
				pongArg, err := strconv.ParseUint(pieces[len(pieces)-1], 10, 64)
				if err == nil {
					client.recordPong(pongArg)
					// set the pong flag, wake if necessary, no-op if set
					select {
					case client.pongEvent <- true:
					default:
					}
				}
			}
		}

		client.lastLine = line
		client.totalLines++
	}
	client.closed <- true
}

// Connect connects to the given server
func (c *Client) Connect(server *Server) error {
	// connect
	var conn net.Conn
	var err error

	addr := strings.TrimPrefix(server.Conn.Address, "unix:")

	if server.Conn.IsTLS {
		conn, err = tls.Dial("tcp", addr, skipVerifyConfig)
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

	go c.readLoop(server)

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
		log.Println(c.Nick, "disconnecting")
		c.SetCloseExpected(true)
		c.Socket.WriteLine("QUIT")
		<-c.closed
	}
}
