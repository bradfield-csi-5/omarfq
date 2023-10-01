package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IPHeader struct {
	Version       uint8
	IHL           uint8
	TotalLength   uint16
	Protocol      uint8
	SourceIP      [4]byte
	DestinationIP [4]byte
}

func parseIPHeader(data []byte) (*IPHeader, error) {
	buffer := make([]byte, 60)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Error reading IP header: %v", err)
	}

	if n < 20 {
		return nil, fmt.Errorf("Data too small to be an IP header")
	}

	header := &IPHeader{
		Version:       data[0] >> 4,
		IHL:           data[0] & 0x0F,
		TotalLength:   binary.BigEndian.Uint16(data[2:4]),
		Protocol:      data[9],
		SourceIP:      [4]byte{data[12], data[13], data[14], data[15]},
		DestinationIP: [4]byte{data[16], data[17], data[18], data[19]},
	}

	return header, nil
}
