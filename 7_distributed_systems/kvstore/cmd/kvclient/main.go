package kvclient

import (
	"fmt"
	"net"
	"strings"

	"github.com/chzyer/readline"
)

const SOCKET_PATH = "/tmp/kvstore.sock"

func main() {
	rl, err := readline.New("Commands: get [key] | set [key]=[value] > ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	instructions := map[string]bool{"set": true, "get": true}

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("Error reading line:", err)
			break
		}

		cmd := strings.Split(line, " ")
		if len(cmd) != 2 {
			fmt.Println("Error: Invalid input. Please use the format 'command [key]' or 'command [key]=[value]'")
			continue
		}

		conn, err := net.Dial("unix", SOCKET_PATH)
		if err != nil {
			fmt.Printf("Failed to connect to server: %s\n", err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Failed to send data: %s\n", err)
			conn.Close()
			continue
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Printf("Failed to read response: %s\n", err)
			conn.Close()
			continue
		}

		fmt.Println("Response from server:", string(response[:n]))

		conn.Close()
	}
}
