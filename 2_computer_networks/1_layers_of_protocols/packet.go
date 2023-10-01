package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type PacketHeader struct {
	CapturedLen uint32
	OriginalLen uint32
}

func skipGlobalHeader(file *os.File) error {
	_, err := file.Seek(GlobalHeaderSize, 0)
	return err
}

func readPacketHeader(file *os.File) (*PacketHeader, error) {
	headerData := make([]byte, PerPacketHeaderSize)
	_, err := file.Read(headerData)
	if err != nil {
		return nil, err
	}
	return &PacketHeader{
		CapturedLen: binary.LittleEndian.Uint32(headerData[8:12]),
		OriginalLen: binary.LittleEndian.Uint32(headerData[12:16]),
	}, nil
}

func printPacketInfo(pHeader *PacketHeader, eHeader *EthernetHeader) {
	fmt.Printf("Captured Length: %d, Original Length: %d\n", pHeader.CapturedLen, pHeader.OriginalLen)
	fmt.Printf("Destination MAC: %x:%x:%x:%x:%x:%x\n", eHeader.DestMAC[0], eHeader.DestMAC[1], eHeader.DestMAC[2], eHeader.DestMAC[3], eHeader.DestMAC[4], eHeader.DestMAC[5])
	fmt.Printf("Source MAC: %x:%x:%x:%x:%x:%x\n", eHeader.SourceMAC[0], eHeader.SourceMAC[1], eHeader.SourceMAC[2], eHeader.SourceMAC[3], eHeader.SourceMAC[4], eHeader.SourceMAC[5])
	switch eHeader.EtherType {
	case IPv4EtherType:
		fmt.Println("Payload is IPv4")
	case IPv6EtherType:
		fmt.Println("Payload is IPv6")
	default:
		fmt.Printf("Unknown EtherType: 0x%x\n", eHeader.EtherType)
	}
}
