package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, err := ln.Accept()
		if err != nil{handleError(err)}
		conns <- conn
	}
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewReader(client)
	for {

		msg, err := reader.ReadString('\n')

		if err != nil {
			msgs <- Message{sender: clientid, message: "LOST CONNECTION TO SERVER"}
			fmt.Println(fmt.Sprintf("[Client%d] lost connection to server Err = %s", clientid, err))
			client.Close()
			return

		} else{
			tidyMsg := fmt.Sprintf("[Client%d]:%s", clientid, msg)
			msgs <- Message{sender: clientid, message: tidyMsg}
			fmt.Printf(tidyMsg)
		}

	}
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	//in order for user to input there own port, they can type -port=PORTNUMBER (e.g. -port=8045)
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()
	ln, _ := net.Listen("tcp", *portPtr)
	//TODO Create a Listener for TCP connections on the port given above.

	conns := make(chan net.Conn)//Create a channel for connections
	msgs := make(chan Message)//Create a channel for messages
	clients := make(map[int]net.Conn)//Create a mapping of IDs to connections

	for {
		go acceptConns(ln, conns)//Start accepting connections

		select {
		case conn := <-conns: //If new connection detected, add to client map with ID and channel
			newClientID := len(clients)
			clients[newClientID] = conn
			//msgs <- Message{-1, fmt.Sprintf("New client, hello [Client%d]", newClientID)}
			//fmt.Print(fmt.Sprintf("New Client: [Client%d]", newClientID))
			go handleClient(clients[newClientID], newClientID, msgs)

			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
		case msg := <-msgs:
			for i:= 0; i < len(clients); i++{
				if i != msg.sender{fmt.Fprintf(clients[i], msg.message)}
				//TODO Deal with a new message
				// Send the message to all clients that aren't the sender
		}

		}
	}
}

