package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type IPHeader struct {
	Version       uint8
	IHL           uint8
	TotalLength   uint16
	Protocol      uint8
	SourceIP      [4]byte
	DestinationIP [4]byte
}

func parseIPHeader(file *os.File) (*IPHeader, error) {
	buffer := make([]byte, IPHeaderMaxSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Error reading IP header: %v", err)
	}

	if n < 20 {
		return nil, fmt.Errorf("Data too small to be an IP header")
	}

	header := &IPHeader{
		Version:       buffer[0] >> 4,
		IHL:           buffer[0] & 0x0F,
		TotalLength:   binary.BigEndian.Uint16(buffer[2:4]),
		Protocol:      buffer[9],
		SourceIP:      [4]byte{buffer[12], buffer[13], buffer[14], buffer[15]},
		DestinationIP: [4]byte{buffer[16], buffer[17], buffer[18], buffer[19]},
	}

	return header, nil
}
