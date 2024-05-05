package utils

import (
	"fmt"
	"strings"
)

func ParseInput(input string) (string, string, error) {
	instructions := map[string]bool{"set": true, "get": true}

	cmd := strings.Split(input, " ")
	if len(cmd) != 2 {
		return "", "", fmt.Errorf("Error: Invalid input -> %s", cmd)
	}

	instruction, item := cmd[0], cmd[1]
	if _, ok := instructions[instruction]; !ok {
		return "", "", fmt.Errorf("Error: Invalid instruction. Please make sure to use either 'get' or 'set'.")
	}

	return instruction, item, nil
}
