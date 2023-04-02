package cligpt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"cligpt/db"
	"cligpt/types"
)

const clearScreen string = "\033[H\033[2J"
const responseColor string = "\x1b[%dm%s\x1b[0m"

type Chunk struct {
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Delta        struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type appEnv struct {
	messages       []types.Message
	model          string
	token          string
	OutputJSON     bool
	isSinglePrompt bool
	InitialPrompt  string
	temperature    float64
	personality    string
	listSessions   bool
	sessions       []types.Session
	currentSession types.Session
}

func (app *appEnv) loadConfig() {
	config := parseConfig()

	app.token = config.Token

	if config.Model == "" {
		app.model = models["chatgpt"]
	} else {
		app.model = config.Model
	}

	app.personality = config.Personality
	app.temperature = config.Temperature
}

func getUserInput() string {
	fmt.Print("> ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func createMessage(role string, content string) types.Message {
	return types.Message{Role: role, Content: content}
}

func printResponse(responseString string) {
	fmt.Print(fmt.Sprintf(responseColor, 32, responseString))
}

func regExpChunk(line []byte) []string {
	pat := regexp.MustCompile(`(data: )(.*)`)
	return pat.FindStringSubmatch(string(line))
}

func parseMessageChunks(resp *http.Response) string {
	var content string

	if resp.StatusCode != 200 {
		log.Fatal(stringifyResponseBody(resp))
	}

	reader := bufio.NewReader(resp.Body)
out:
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			break
		}

		var chunk Chunk

		matches := regExpChunk(line)

		if len(matches) > 1 {
			if err := json.Unmarshal([]byte(strings.Trim(matches[2], " ")), &chunk); err != nil {
				log.Fatal("Error parsing response body:", err)
			}
			if chunk.Choices[0].Delta.Content != "" {
				fmt.Print(fmt.Sprintf(responseColor, 32, chunk.Choices[0].Delta.Content))
				content += chunk.Choices[0].Delta.Content
			}

			if chunk.Choices[0].FinishReason == "stop" {
				break out
			}
		}
	}

	return content
}

func (app *appEnv) singlePrompt() {
	fmt.Print(clearScreen)
	req := buildRequest(app)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	if app.OutputJSON {
		printResponse(stringifyResponseBody(resp))
		print()
		return
	}

	responseBody := parseResponse(resp)

	printResponse(responseBody.Choices[0].Message.Content)
	print()
}

func (app *appEnv) sessionPrompt() {
	fmt.Print(clearScreen)

	req := buildRequest(app)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	content := parseMessageChunks(resp)

	app.messages = append(app.messages, types.Message{Role: "assistant", Content: content})

	if reflect.ValueOf(app.currentSession).IsZero() {
		app.currentSession = db.CreateSession(app.messages)
	} else {
		db.UpdateSession(app.currentSession.ID, app.messages)
	}

	if app.currentSession.ID != 0 {
	} else {
	}
	fmt.Println()
}

func Init() {
	createConfig()
	db.InitDB()
	SelectAndSaveModel()
	GetAndSaveToken()
}

func InitApp() appEnv {
	app := appEnv{}
	app.loadConfig()
	return app
}

func (app *appEnv) Chat() {
	app.loadConfig()
	if app.personality != "" {
		app.messages = append(app.messages, createMessage("system", app.personality))
	}
	for true {
		var input string
		if app.InitialPrompt != "" {
			input = app.InitialPrompt
			app.InitialPrompt = ""
		} else {
			input = getUserInput()
		}
		app.messages = append(app.messages, createMessage("user", input))
		app.sessionPrompt()
	}
}

func (app *appEnv) SinglePrompt() {
	app.loadConfig()
	app.isSinglePrompt = true
	app.messages = append(app.messages, createMessage("user", app.InitialPrompt))
	app.singlePrompt()
}

func (app *appEnv) ListAndSelectSession() {
	app.loadConfig()
	app.listSessions = true
	app.sessions = db.GetLastTenSessions()

	sessionNames := []string{}
	for _, e := range app.sessions {
		sessionNames = append(sessionNames, e.Messages[0].Content)
	}

	selectSessionPromptContent := promptSelectContent{
		label:        "Select a previous chat",
		selectValues: sessionNames,
	}
	promptResult := promptGetSelect(selectSessionPromptContent)

	app.currentSession = app.sessions[promptResult.index]
	app.messages = app.currentSession.Messages
	for _, e := range app.currentSession.Messages {
		if e.Role == "user" {
			fmt.Println(e.Content)
		} else {
			printResponse(e.Content)
			fmt.Println()
		}
	}
}
