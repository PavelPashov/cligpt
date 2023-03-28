package main

import (
	"flag"
	"fmt"
)

var session = false

func continueAfterHandlingArgs(args []string) bool {
	var help string
	var token string
	var model string

	flag.StringVar(&help, "h", "help", "Usage: cligpt <prompt>")
	flag.StringVar(&token, "t", "token", "Provide a token used for authentication.")
	flag.StringVar(&model, "m", "model", "Please provide a model.")
	flag.BoolVar(&session, "s", false, "Provide flag if you want to start a session")
	flag.Parse()

	if isFlagPassed("t") {
		config := parseConfig()
		config.Token = token
		saveConfig(config)
		fmt.Println("Token saved to config file at: ", getConfigPath())
		return false
	}

	if isFlagPassed("h") {
		flag.PrintDefaults()
		return false
	}

	if isFlagPassed("m") {
		if model != ChatGpt && model != Gpt4 {
			fmt.Println("Please provide a valid model.")
			return false
		}
		config := parseConfig()
		config.Model = model
		saveConfig(config)
		fmt.Println("Model saved to config file at: ", getConfigPath())
		return false
	}

	if isFlagPassed("s") {
		session = true
	}

	if len(args) == 0 {
		fmt.Println("Please provide at least one argument.")
		return false
	}

	return true
}
