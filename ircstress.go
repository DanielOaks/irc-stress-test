// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/DanielOaks/irc-stress-test/stress"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `ircstress.
Usage:
	ircstress run [--clients=<num>] [--chan-ratio=<ratio>] [--chan-join-percent=<ratio>] [--queues=<num>] [--wait] <server-details>...
	ircstress -h | --help
	ircstress --version
Options:
	--clients=<num>              The number of clients that should connect [default: 10000].
	--channels=<chans>           How many channels exist, limited to number of clients [default: 1000].
	--chan-join-percent=<ratio>  How likely each client is to join one or more channels [default: 0.9].

	--queues=<num>     How many queues to run events on, limited to number of clients [default: 3].
	--wait             After each action, waits for server response before continuing.
	<server-details>   Set of server details, of the format: "Name,Addr,TLS", where Addr is like "localhost:6667" and TLS is either "yes" or "no".

	-h --help          Show this screen.
	--version          Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, stress.SemVer, false)

	if arguments["run"].(bool) {
		// run string
		var optionString string
		if !arguments["--wait"].(bool) {
			optionString += "not "
		}
		optionString += "waiting"

		fmt.Println(fmt.Sprintf("Running tests (%s)", optionString))

		// assemble each server's details
		servers := make(map[string]*stress.Server)
		for _, serverString := range arguments["<server-details>"].([]string) {
			serverList := strings.Split(serverString, ",")
			if len(serverList) != 3 {
				log.Fatal("Could not parse server details string:", serverString)
			}

			var isTLS bool
			if strings.ToLower(serverList[2]) == "yes" {
				isTLS = true
			} else if strings.ToLower(serverList[2]) == "no" {
				isTLS = false
			} else {
				log.Fatal("TLS must be either 'yes' or 'no', could not parse whether to enable TLS from server details:", serverString)
			}

			newServer := stress.Server{
				Name: serverList[0],
				Conn: stress.ServerConnectionDetails{
					Address: serverList[1],
					IsTLS:   isTLS,
				},
			}

			fmt.Println("Running server", newServer.Name, ":", newServer.Conn.Address)

			servers[newServer.Name] = &newServer
		}

		// create event queues
		eventQueues := make([]stress.EventQueue, 0)

		// make clients
		clientCount, err := strconv.Atoi(arguments["--clients"].(string))
		if err != nil || clientCount < 1 {
			log.Fatal("Not a real number of clients:", arguments["--clients"].(string))
		}

		clients := make(map[int]*stress.Client)
		for i := 0; i < clientCount; i++ {
			var newClient *stress.Client
			newClient = &stress.Client{
				Nick: fmt.Sprintf("cli%d", i),
			}

			clients[i] = newClient

			// for now we'll just have one event list per client for simplicity
			events := stress.NewEventQueue()
			events.Events = append(events.Events, stress.Event{
				Client: i,
				Type:   stress.ETConnect,
			})

			// send NICK+USER
			// events.Events = append(events.Events, stress.Event{
			// 	Client: i,
			// 	Type:   stress.ETLine,
			// 	Line:   fmt.Sprintf("CAP END\r\n", newClient.Nick),
			// })
			events.Events = append(events.Events, stress.Event{
				Client: i,
				Type:   stress.ETLine,
				Line:   fmt.Sprintf("NICK %s\r\n", newClient.Nick),
			})
			events.Events = append(events.Events, stress.Event{
				Client: i,
				Type:   stress.ETLine,
				Line:   "USER test 0 * :I am a cool person!\r\n",
			})

			//TODO(dan): send NICK/USER
			events.Events = append(events.Events, stress.Event{
				Client: i,
				Type:   stress.ETDisconnect,
			})

			eventQueues = append(eventQueues, events)
		}

		// run for each server
		for name, server := range servers {
			fmt.Println("Testing", name)

			// start each event queue
			for _, events := range eventQueues {
				go events.Run(server, clients)
			}

			// wait for each of them to be finished
			for _, events := range eventQueues {
				<-events.Finished
			}
		}
	}
}
