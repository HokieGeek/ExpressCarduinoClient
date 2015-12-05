package main

import (
	"fmt"
	"github.com/hokiegeek/ExpressCarduinoDaemon/connection"
	"log"
	"syscall"
)

func main() {
	// TODO: store historical values?

	// Connect to the board
	// s := &connection.Serial{DeviceName: "/dev/expresscarduino", BaudRate: syscall.B115200}
	s := &connection.Serial{DeviceName: "/dev/arduinoMetro328", BaudRate: syscall.B9600}
	conn, err := connection.New(s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", conn)

	err = conn.Connect()
	if err != nil {
		// TODO: if no connection, keep trying periodically
		log.Fatal(err)
	}

	fmt.Printf("%s\n", conn)

	// TODO: Kick off a routine that just reads bytes
}
