package main

import (
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
)

const FILENAME = "kvstore.json"

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
			if err := writeToFile(FILENAME, key, val); err != nil {
				fmt.Printf("Error: Could not write to JSON file. %s\n", err)
				continue
			}
			fmt.Printf("Inserted key: %s and value: %s\n", key, val)
		case "get":
			val, err := getFromFile(FILENAME, value)
			if err != nil {
				fmt.Printf("Error: Could not read from JSON file. %s\n", err)
				continue
			}
			fmt.Println(val)
		}
	}
}

func writeToFile(filename, key, value string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		if err != io.EOF {
			return err
		}
		data = make(map[string]string)
	}

	data[key] = value

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func getFromFile(filename, key string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("Key %q not found in the file", key)
	}

	return val, nil
}
