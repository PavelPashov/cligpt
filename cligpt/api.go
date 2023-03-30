package cligpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	URL string = "https://api.openai.com/v1/chat/completions"
)

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	}
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func buildRequest(app *appEnv) *http.Request {
	var reqBody RequestBody

	reqBody.Model = app.model
	reqBody.Stream = !app.isSinglePrompt

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
