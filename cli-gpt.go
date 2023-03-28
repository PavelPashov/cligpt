package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	dev              string = "developer"
	devMotivation    string = "You provide short and concise answers about development in markdown."
	writer           string = "writer"
	writerMotivation string = "You are a writer who is creative and imaginative, you provide long and detailed answers."
)

const (
	URL string = "https://api.openai.com/v1/chat/completions"
)

const (
	ChatGpt string = "gpt-3.5-turbo"
	Gpt4    string = "gpt-4"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	}
}

type ReqBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func createRequest(messages []Message, input string, config Config) *http.Request {
	url := "https://api.openai.com/v1/chat/completions"

	var reqBody ReqBody

	reqBody.Model = config.Model

	reqBody.Messages = append(messages, Message{Role: "user", Content: input})

	finalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(finalReqBody))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	return req
}

func parseResponse(resp *http.Response) Response {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var responseBody Response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Fatal("Error parsing response body:", err)
	}

	return responseBody
}

func askQuestion(messages []Message, input string, config Config) {
	req := createRequest(messages, input, config)

	// This shows the loading spinner
	shutdownCh := make(chan struct{})
	go showSpinner(shutdownCh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	close(shutdownCh)

	respBody := parseResponse(resp)

	message := respBody.Choices[0].Message
	messages = append(messages, message)

	fmt.Print(clearScreen)
	fmt.Println(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "Response:"))
	fmt.Println(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, message.Content))
}
