package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/omarfq/kvstore/internal/store"
	"github.com/omarfq/kvstore/internal/utils"
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
		fmt.Println("Server running...")
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

	instruction, item, err := utils.ParseInput(string(data))

	response, err := processInstruction(instruction, item)
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
		return fmt.Sprint("OK\n"), nil
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
