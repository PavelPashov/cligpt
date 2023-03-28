package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	clearScreen string = "\033[H\033[2J"
)

func main() {
	fmt.Print(clearScreen)
	args := os.Args[1:]

	if !continueAfterHandlingArgs(args) {
		return
	}

	config := parseConfig()

	if config.Token == "" {
		fmt.Println("Please provide a valid token first.")
		return
	}

	if session {
		var messages []Message
		for true {
			input := getUserInput()
			messages = askQuestion(messages, input, config)
		}
	} else {
		askQuestion(nil, strings.Join(args[:], " "), config)
	}

}
