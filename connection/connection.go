package connection

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

const initialBaudRate uint32 = syscall.B9600
const connectionTimeout time.Duration = time.Second * 10
const handshakeRequestChar byte = ','
const disconnectRequestChar byte = '.'
const AckChar byte = ':'

type ConnectionState int

const (
	Inactive ConnectionState = iota
	Handshaking
	Active
)

func (s ConnectionState) String() string {
	switch s {
	case Inactive:
		return "Inactive"
	case Handshaking:
		return "Handshaking"
	case Active:
		return "Active"
	}
	return "Unknown"
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
	return uint8(vtime)
}

type Serial struct {
	DeviceName string
	BaudRate   uint32
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
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(c.file.Fd()),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&term)),
	); errno != 0 {
		// TODO: include errno in this
		return errors.New("Encountered error doing IOCTL syscall")
	}

	if err := syscall.SetNonblock(int(c.file.Fd()), false); err != nil {
		return err
	}

	return nil
}

func (c *Connection) performHandshake() error {
	c.State = Handshaking

	_, err := c.Write([]byte{handshakeRequestChar})
	// _, err = c.Write([]byte{handshakeRequestChar, c.serial.BaudRate})
	if err != nil {
		return err
	}

	// Wait for AckChar
	const MAX_ATTEMPTS = 40
	buf := make([]byte, 1)
	var n int
	var attempts int
	for n <= 0 && attempts < MAX_ATTEMPTS {
		n, err = c.Read(buf)
		if err != nil {
			return err
		}
		attempts++
		time.Sleep(100 * time.Millisecond)
	}
	if buf[0] != AckChar {
		return errors.New("Did not receive acknowledgment. Cannot connect")
	}

	return nil
}

func (c *Connection) Connect() error {
	// Open the file
	var err error
	c.file, err = os.OpenFile(c.serial.DeviceName, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return err
	}

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
	   err = c.setBaudRate(c.serial.BaudRate)
	   if err != nil {
	       return err
	   }
	*/

	c.State = Active

	return nil
}

func (c *Connection) Disconnect() (err error) {
	_, err = c.Write([]byte{disconnectRequestChar})
	// _, err = c.Write([]byte{handshakeRequestChar, c.serial.BaudRate})
	if err != nil {
		return err
	}
	return c.file.Close()
}

func (c *Connection) Read(b []byte) (n int, err error) {
	return c.file.Read(b)
}

func (c *Connection) Write(b []byte) (n int, err error) {
	return c.file.Write(b)
}

func (c *Connection) String() string {
	var buf bytes.Buffer

	// Device name
	buf.WriteString("Device: ")
	buf.WriteString(c.serial.DeviceName)
	buf.WriteString("\n")

	// Baud rate
	buf.WriteString("Baud rate: ")
	buf.WriteString("TODO")
	buf.WriteString(strconv.Itoa(int(c.serial.BaudRate))) // TODO: whoops
	buf.WriteString("\n")

	// State
	buf.WriteString("Connection state: ")
	buf.WriteString(c.State.String())
	buf.WriteString("\n")

	return buf.String()
}

func New(ser *Serial) (*Connection, error) {
	c := new(Connection)
	c.serial = ser
	return c, nil
}
