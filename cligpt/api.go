package cligpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/eitamonya/cligpt/types"
)

const (
	CHAT_URL  string = "https://api.openai.com/v1/chat/completions"
	IMAGE_URL string = "https://api.openai.com/v1/images/generations"
)

type ChatResponseBody struct {
	Choices []struct {
		Message types.Message `json:"message"`
	}
}

type ChatRequestBody struct {
	Model       string          `json:"model"`
	Messages    []types.Message `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type ImageResponseBody struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type ImageRequestBody struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func buildCompletionRequest(app *appEnv) *http.Request {
	var reqBody ChatRequestBody

	reqBody.Model = app.model
	reqBody.Stream = !app.isSinglePrompt
	reqBody.Temperature = app.temperature
	reqBody.MaxTokens = app.max_tokens
	reqBody.Messages = app.currentSession.Messages

	finalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	req, err := http.NewRequest("POST", CHAT_URL, bytes.NewBuffer(finalReqBody))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+app.token)

	return req
}

func parseCompletionResponse(resp *http.Response) ChatResponseBody {
	if strings.HasPrefix(http.StatusText(resp.StatusCode), "4") || strings.HasPrefix(http.StatusText(resp.StatusCode), "5") {
		log.Fatal(stringifyResponseBody(resp))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var responseBody ChatResponseBody

	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Fatal("Error parsing response body:", err)
	}

	return responseBody
}

func buildImageRequest(app *appEnv) *http.Request {
	var reqBody ImageRequestBody

	reqBody.Prompt = app.InitialPrompt
	reqBody.N = 1
	reqBody.Size = "1024x1024"

	finalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	req, err := http.NewRequest("POST", IMAGE_URL, bytes.NewBuffer(finalReqBody))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+app.token)

	return req
}

func parseImageResponse(resp *http.Response) ImageResponseBody {
	if strings.HasPrefix(http.StatusText(resp.StatusCode), "4") || strings.HasPrefix(http.StatusText(resp.StatusCode), "5") {
		log.Fatal(stringifyResponseBody(resp))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var responseBody ImageResponseBody

	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Fatal("Error parsing response body:", err)
	}

	return responseBody
}

func stringifyResponseBody(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	finalBody := &bytes.Buffer{}
	if err := json.Indent(finalBody, body, "", "  "); err != nil {
		panic(err)
	}

	return finalBody.String()
}
