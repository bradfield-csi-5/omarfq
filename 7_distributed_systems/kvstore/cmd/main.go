package main

import (
	"fmt"
	"github.com/chzyer/readline"
)

func main() {
	rl, err := readline.New("Wirte something here > ")

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		fmt.Println(line)
	}
}
