package main

import (
	"bytes"
	"encoding/binary"
	"github.com/hokiegeek/ExpressCarduinoDaemon/connection"
	"log"
	"syscall"
)

func connect(device string) (*connection.Connection, error) {
	// Connect to the board
	// s := &connection.Serial{DeviceName: device, BaudRate: syscall.B115200}
	s := &connection.Serial{DeviceName: device, BaudRate: syscall.B9600}
	conn, err := connection.New(s)
	if err != nil {
		return nil, err
	}

	defer func() {
		conn.Disconnect()
	}()

	log.Printf("Establishing connection to device:\n%s\n", conn)
	err = conn.Connect()
	if err != nil {
		// TODO: if no connection, keep trying periodically
		return nil, err
	}
	log.Printf("Connected to device:\n%s\n", conn)

	return conn, nil
}

func toggleData(conn *connection.Connection) error {
	// TODO: store historical values?
	_, err := conn.Write([]byte("T"))
	if err != nil {
		return err // TODO
	}

	return nil
}

// func main() {
// log.Printf("Int: %

func main() {
	/*
		// conn, err := connect("/dev/expresscarduino")
		conn, err := connect("/dev/arduinoMetro328")
		if err != nil {
			panic(err) // TODO
		}

		toggleData(conn) // TODO: specify type of data to toggle
	*/

	// s := &connection.Serial{DeviceName: device, BaudRate: syscall.B115200}
	// s := &connection.Serial{DeviceName: device, BaudRate: syscall.B9600}
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

	// Toggle button data
	_, err = conn.Write([]byte("B"))
	if err != nil {
		panic(err)
	}

	// TODO: Kick off a routine that just reads bytes
	log.Printf("Reading stream...\n")
	cmd := make([]byte, 1)
	for conn.State == connection.Active {
		n, err := conn.Read(cmd)
		if err != nil {
			log.Fatal(err)
		}

		if string(cmd[:n]) == "B" {
			buf := make([]byte, 2)
			n, err := conn.Read(buf)
			log.Printf("Got %d bytes\n", n)
			if err != nil {
				log.Fatal(err)
			}
			if n != 2 {
				log.Fatal("Did not read correct number of bytes!")
			}
			
			var val int16
			n, err = binary.LittleEndian.PutUint16(buf, uint16(val)) // Arduino is little endian, like most uCs
			if err != nil {
				log.Fatal(err)
			}
			
			log.Print("Button: %d", val)
		}

		// log.Printf("(%d) %q", n, buf[:n])
	}
	log.Printf("Ended connection to device:\n%s\n", conn)
	// TODO: Attempt to reconnect
}
