package cligpt

import (
	"bytes"
	"cligpt/types"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	URL string = "https://api.openai.com/v1/chat/completions"
)

type ResponseBody struct {
	Choices []struct {
		Message types.Message `json:"message"`
	}
}

type RequestBody struct {
	Model       string          `json:"model"`
	Messages    []types.Message `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

func buildRequest(app *appEnv) *http.Request {
	var reqBody RequestBody

	reqBody.Model = app.model
	reqBody.Stream = !app.isSinglePrompt

	if app.temperature != 0 {
		reqBody.Temperature = app.temperature
	}

	if app.max_tokens != 0 {
		reqBody.MaxTokens = app.max_tokens
	}

	reqBody.Messages = app.messages

	finalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Error creating request body:", err)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(finalReqBody))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+app.token)

	return req
}

func parseResponse(resp *http.Response) ResponseBody {
	if strings.HasPrefix(http.StatusText(resp.StatusCode), "4") || strings.HasPrefix(http.StatusText(resp.StatusCode), "5") {
		log.Fatal(stringifyResponseBody(resp))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	var responseBody ResponseBody

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
