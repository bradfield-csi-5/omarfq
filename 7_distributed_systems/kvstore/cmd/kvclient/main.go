package kvclient

import (
	"fmt"
	"net"

	"github.com/chzyer/readline"
	"github.com/omarfq/kvstore/internal/utils"
)

const SOCKET_PATH = "/tmp/kvstore.sock"

func main() {
	rl, err := readline.New("Commands: get [key] | set [key]=[value] > ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("Error reading line:", err)
			break
		}

		_, _, err = utils.ParseInput(line)
		if err != nil {
			fmt.Printf("Failed to parse input: %s", err)
			continue
		}

		// Connect to the server via UNIX socket
		conn, err := net.Dial("unix", SOCKET_PATH)
		if err != nil {
			fmt.Printf("Failed to connect to server: %s\n", err)
			continue
		}

		// Send the instruction and value to the server
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Failed to send data: %s\n", err)
			conn.Close()
			continue
		}

		// Read response from server
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Printf("Failed to read response: %s\n", err)
			conn.Close()
			continue
		}

		fmt.Println("Response from server:", string(response[:n]))

		// Close the connection after the interaction
		conn.Close()
	}
}
