package main

import (
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

	defer func() {
		conn.Disconnect()
	}()

	log.Printf("Establishing connection to device:\n%s\n", conn)
	err = conn.Connect()
	if err != nil {
		// TODO: if no connection, keep trying periodically
		panic(err)
	}
	log.Printf("Connected to device:\n%s\n", conn)

	_, err = conn.Write([]byte("T"))
	if err != nil {
		log.Fatal("sigh")
	}

	// TODO: Kick off a routine that just reads bytes
	log.Printf("Reading stream...\n")
	buf := make([]byte, 128)
	for conn.State == connection.Active {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if string(buf[:n]) == "T" {
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%d", buf[:n])
		}

		// log.Printf("(%d) %q", n, buf[:n])
	}
	log.Printf("Ended connection to device:\n%s\n", conn)
	// TODO: Attempt to reconnect
}
