package kvserver

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/omarfq/kvstore/pkg/store"
)

const SOCKET_PATH = "/tmp/kvstore.sock"
const PATH = "data/kvstore.json"

func main() {
	if err := os.Remove(SOCKET_PATH); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	listener, err := net.Listen("unix", SOCKET_PATH)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %s\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	readData(conn)
}

func readData(conn net.Conn) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading data: %s\n", err)
		return
	}
	data := string(buf[:n])
	parts := strings.SplitN(data, " ", 2)
	if len(parts) != 2 {
		fmt.Fprintf(conn, "Error: Invalid input\n")
		return
	}

	response, err := processInstruction(parts[0], parts[1])
	if err != nil {
		fmt.Fprintf(conn, "Error: %s\n", err)
		return
	}
	fmt.Fprintf(conn, "%s\n", response)
}

func processInstruction(instruction, item string) (string, error) {
	kvstore := store.NewFileKVStore(PATH)

	switch instruction {
	case "set":
		itemSlice := strings.Split(item, "=")
		if len(itemSlice) != 2 {
			return "", fmt.Errorf("invalid key-item pair")
		}
		key, val := itemSlice[0], itemSlice[1]
		if err := kvstore.Set(key, val); err != nil {
			return "", fmt.Errorf("could not write to JSON file: %s", err)
		}
		return fmt.Sprintf("Inserted key: %s and item: %s", key, val), nil
	case "get":
		val, err := kvstore.Get(item)
		if err != nil {
			return "", fmt.Errorf("could not read from JSON file: %s", err)
		}
		return val, nil
	default:
		return "", fmt.Errorf("unknown instruction")
	}
}
