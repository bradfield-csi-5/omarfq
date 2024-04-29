package kvclient

import (
	"fmt"
	"github.com/chzyer/readline"
	"strings"

	"github.com/omarfq/kvstore/pkg/store"
)

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

	}
}
