package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
}

type ReqBody struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

func saveToken(token string) string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(homedir, ".cligpt")
	_ = os.Mkdir(path, os.ModePerm)

	f, err := os.Create(filepath.Join(path, "token"))
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte(token))

	defer f.Close()
	return f.Name()
}

func getToken() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(homedir, ".cligpt", "token")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Please run cligpt -t <token> to save your token.")
		log.Fatal(err)
	}

	return string(file)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	fmt.Print(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "Loading"))
	for {
		select {
		case <-ticker.C:
			fmt.Print(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "."))
		case <-shutdownCh:
			return
		}
	}
}

func main() {
	fmt.Print("\033[H\033[2J")
	var help string
	var token string
	args := os.Args[1:]

	flag.StringVar(&help, "h", "help", "Usage: cligpt <prompt>")
	flag.StringVar(&token, "t", "token", "Provide a token used for authentication.")
	flag.Parse()

	if isFlagPassed("t") {
		fmt.Println("Token saved to", saveToken(token))
		return
	}

	if isFlagPassed("h") {
		flag.PrintDefaults()
		return
	}

	accessToken := getToken()

	if accessToken == "" {
		fmt.Println("Please provide a valid token first.")
		return
	}

	if len(args) == 0 {
		fmt.Println("Please provide at least one argument.")
		return
	}

	url := "https://api.openai.com/v1/chat/completions"

	data := []byte(`{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": ""}]}`)

	var reqBody ReqBody

	if err := json.Unmarshal(data, &reqBody); err != nil {
		fmt.Println("Error parsing request body:", err)
		return
	}

	reqBody.Messages[0].Content = strings.Join(args[:], " ")

	finalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(finalReqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	shutdownCh := make(chan struct{})
	go indicator(shutdownCh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	close(shutdownCh)

	fmt.Print("\033[H\033[2J")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println(string(body))

	var responseBody Response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		fmt.Println("Error parsing response body:", err)
		return
	}

	value := responseBody.Choices[0].Message.Content
	fmt.Println(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "Response:"))
	fmt.Println(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, value))
}
