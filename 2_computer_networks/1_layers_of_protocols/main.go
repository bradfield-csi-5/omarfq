package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("/Users/omarflores/Documents/BradfieldCSI/omarfq/2_computer_networks/1_layers_of_protocols/net.cap")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	err = skipGlobalHeader(file)
	if err != nil {
		fmt.Println("Error skipping global header:", err)
		return
	}

	for {
		pHeader, err := readPacketHeader(file)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error reading packet header:", err)
			return
		}

		eHeader, err := readEthernetHeader(file)
		if err != nil {
			fmt.Println("Error reading Ethernet header:", err)
			return
		}

		ipHeader, err := parseIPHeader(file)

		printPacketInfo(pHeader, eHeader, ipHeader)

		remainingPacketData := pHeader.CapturedLen
		remainingPacketData -= EthernetHeaderSize
		remainingPacketData -= IPHeaderMaxSize

		// Skip remaining packet data
		file.Seek(int64(remainingPacketData), 1)
	}
}

func printPacketInfo(pHeader *PacketHeader, eHeader *EthernetHeader, ipHeader *IPHeader) {
	// Packet Header Info
	fmt.Printf("Captured Length: %d, Original Length: %d\n", pHeader.CapturedLen, pHeader.OriginalLen)

	// Ethernet Header Info
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

	// IP Header info
	fmt.Printf("Source IP: %v.%v.%v.%v\n", ipHeader.SourceIP[0], ipHeader.SourceIP[1], ipHeader.SourceIP[2], ipHeader.SourceIP[3])
	fmt.Printf("Destination IP: %v.%v.%v.%v\n", ipHeader.DestinationIP[0], ipHeader.DestinationIP[1], ipHeader.DestinationIP[2], ipHeader.DestinationIP[3])
	fmt.Printf("Protocol: %v\n", ipHeader.Protocol)

	// Compute payload length by subtracting header length from total length
	payloadLength := ipHeader.TotalLength - (uint16(ipHeader.IHL) * 4)
	fmt.Printf("Payload Length: %d bytes\n", payloadLength)

	fmt.Println("========================")
}
