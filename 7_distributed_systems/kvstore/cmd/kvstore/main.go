package main

import (
	"fmt"
	"github.com/chzyer/readline"
	"strings"

	"github.com/omarfq/kvstore/pkg/store"
)

const PATH = "data/kvstore.json"

func main() {
	rl, err := readline.New("Commands: get [key] | set [key]=[value] > ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	instructions := map[string]bool{"set": true, "get": true}

	kvstore := store.NewFileKVStore(PATH)

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		cmd := strings.Split(line, " ")
		if len(cmd) != 2 {
			fmt.Println("Error: Invalid input.")
			continue
		}

		instruction, value := cmd[0], cmd[1]
		if _, ok := instructions[instruction]; !ok {
			fmt.Println("Error: Invalid instruction. Please make sure to use either 'get' or 'set'.")
			continue
		}

		switch instruction {
		case "set":
			valueSlice := strings.Split(value, "=")
			if len(valueSlice) != 2 {
				fmt.Println("Error: Invalid key-value pair")
				continue
			}
			key, val := valueSlice[0], valueSlice[1]
			if err := kvstore.Set(key, val); err != nil {
				fmt.Printf("Error: Could not write to JSON file. %s\n", err)
				continue
			}
			fmt.Printf("Inserted key: %s and value: %s\n", key, val)
		case "get":
			val, err := kvstore.Get(value)
			if err != nil {
				fmt.Printf("Error: Could not read from JSON file. %s\n", err)
				continue
			}
			fmt.Println(val)
		}
	}
}
