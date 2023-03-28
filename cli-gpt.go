package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	Stream   bool      `json:"stream"`
}

type Chunk struct {
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Delta        struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func createRequest(messages []Message, input string, config Config) *http.Request {
	url := "https://api.openai.com/v1/chat/completions"

	var reqBody ReqBody

	reqBody.Model = config.Model
	reqBody.Stream = true

	reqBody.Messages = messages

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

func createMessage(role string, content string) Message {
	return Message{Role: role, Content: content}
}

func askQuestion(messages []Message, input string, config Config) []Message {
	fmt.Print(clearScreen)

	newMessage := createMessage("user", input)
	messages = append(messages, newMessage)

	req := createRequest(messages, input, config)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	var content string

	reader := bufio.NewReader(resp.Body)
out:
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		pat := regexp.MustCompile(`(data: )(.*)`)
		match := pat.FindStringSubmatch(string(line))

		var chunk Chunk

		if len(match) > 1 {
			if err := json.Unmarshal([]byte(strings.Trim(match[2], " ")), &chunk); err != nil {
				log.Fatal("Error parsing response body:", err)
			}
			if chunk.Choices[0].Delta.Content != "" {
				fmt.Print(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, chunk.Choices[0].Delta.Content))
				content += chunk.Choices[0].Delta.Content
			}

			if chunk.Choices[0].FinishReason == "stop" {
				break out
			}
		}
	}

	messages = append(messages, Message{Role: "assistant", Content: content})
	fmt.Println()
	return messages
}
