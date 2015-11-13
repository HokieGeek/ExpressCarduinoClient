package connection

import (
	// "errors"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const initialBaudRate uint32 = syscall.B9600
const connectionTimeout time.Duration = time.Second * 10

type ConnectionState int

const (
	Inactive ConnectionState = iota
	Handshaking
	Active
)

func (s ConnectionState) String() string {
	if s&Inactive == Inactive {
		return "Inactive"
	}
	return ""
}

func getVtime(duration time.Duration) uint8 {
	const (
		MINTIMEOUT = 1
		MAXTIMEOUT = 255
	)
	vtime := (duration.Nanoseconds() / 1e6 / 100)
	if vtime < MINTIMEOUT {
		vtime = MINTIMEOUT
	} else if vtime > MAXTIMEOUT {
		vtime = MAXTIMEOUT
	}
	return vtime
}

type Serial struct {
	deviceName string
	baudRate   uint32
}

type Connection struct {
	serial *Serial
	file   *os.File
	State  ConnectionState
}

func (c *Connection) setBaudRate(rate uint32) error {
	// Create the term IO settings structure
	term := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | rate,
		Cc:     [32]uint8{syscall.VMIN: 0, syscall.VTIME: getVtime(connectionTimeout)},
		Ispeed: rate,
		Ospeed: rate,
	}
	
	// Make the IOCTL system call to configure the term
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(c.file.Fd()),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&term)),
	)
	if err != nil {
		return err
	}
	
	// Is this necessary?
	/*
	err = syscall.SetNonblock(int(c.file.Fd()), true)
	if err != nil {
		return err
	}
	*/

	return nil
}

func (c *Connection) performHandshake() error {
	// TODO: Look for handshake query byte and return response
	return nil
}

func (c *Connection) Connect() error {
	// Open the file
	var err error
	c.file, err = os.OpenFile(c.serial.deviceName, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return err
	}
	// TODO: I need to close this file

	// Create a connection using a safe baud rate
	err = c.setBaudRate(initialBaudRate)
	if err != nil {
		return err
	}

	// Perform handshake with EC
	err = c.performHandshake()
	if err != nil {
		return err
	}

	// Change baud rate to what was requested
	/*
	   err = c.setBaudRate(c.serial.baudRate)
	   if err != nil {
	       return err
	   }
	*/

	return nil
}

func (c *Connection) String() string {
	// Device name
	// Baud rate
	// State
	return "TODO"
}

func New(ser Serial) (*Connection, error) {
	c := new(Connection)
	c.serial = &ser
	return c, nil
}
