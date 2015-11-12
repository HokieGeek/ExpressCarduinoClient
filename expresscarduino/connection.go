package expresscarduino

import (
    "os"
    "syscall"
    // "time"
    // "unsafe"
)

const initialBaudRate uint32 = syscall.B9600

type Serial struct {
    deviceName string
}

type Connection struct {
    file *os.File
}

func (c *Connection) Connect() (error) {
    // Open the file
    // do an IOCTL call to set the parity and initial baud rate
    // Perform handshake with EC
    return nil
}

func New(ser Serial) (*Connection, error) {
    return nil, nil
}
