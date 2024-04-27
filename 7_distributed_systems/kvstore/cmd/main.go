package main

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("Wirte something here > ")

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	kvStore := map[string]string{}
	instructions := map[string]bool{"set": true, "get": true}

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
			kvStore[key] = val
			fmt.Printf("Inserted key: %s and value: %s\n", key, val)
		case "get":
			if _, ok := kvStore[value]; !ok {
				fmt.Println("Error: The key does not exist in the key-value store.")
				continue
			}
			fmt.Println(kvStore[value])
		}
	}
}
