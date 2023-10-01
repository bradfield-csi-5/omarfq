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

		//ipHeaderData := make([]byte

		printPacketInfo(pHeader, eHeader)

		// Skip remaining packet data
		file.Seek(int64(pHeader.CapturedLen-EthernetHeaderSize), 1)
	}
}
