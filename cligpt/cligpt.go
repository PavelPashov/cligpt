package cligpt

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/eitamonya/cligpt/types"

	"github.com/eitamonya/cligpt/db"
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
	// messages       []types.Message
	model          string
	token          string
	OutputJSON     bool
	isSinglePrompt bool
	InitialPrompt  string
	temperature    float64
	max_tokens     int
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

	var personality string
	for _, p := range config.Personalities {
		if p.Active {
			personality = p.Context
		}
	}

	app.personality = personality
	app.temperature = config.Temperature
	app.max_tokens = config.MaxTokens
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
	req := buildCompletionRequest(app)

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

	responseBody := parseCompletionResponse(resp)

	printResponse(responseBody.Choices[0].Message.Content)
	print()
}

func (app *appEnv) sessionPrompt() {
	fmt.Print(clearScreen)

	req := buildCompletionRequest(app)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	content := parseMessageChunks(resp)

	// app.messages = append(app.messages, types.Message{Role: "assistant", Content: content})
	// newMessages := append(app.currentSession.Messages, types.Message{Role: "assistant", Content: content})

	app.currentSession.Messages = append(app.currentSession.Messages, types.Message{Role: "assistant", Content: content})

	if app.currentSession.ID == 0 {
		app.currentSession = db.CreateSession(app.currentSession.Messages)
	} else {
		db.UpdateSession(app.currentSession.ID, app.currentSession.Messages)
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

	if app.currentSession.ID == 0 {
		app.currentSession = types.Session{Messages: []types.Message{}}
	}

	if app.personality != "" && app.currentSession.ID == 0 {
		app.currentSession.Messages = append(app.currentSession.Messages, createMessage("system", app.personality))
	}

	for true {
		var input string
		if app.InitialPrompt != "" {
			input = app.InitialPrompt
			app.InitialPrompt = ""
		} else {
			input = getUserInput()
		}

		if input == "exit" || input == "quit" || input == "q" {
			break
		}

		app.currentSession.Messages = append(app.currentSession.Messages, createMessage("user", input))
		app.sessionPrompt()
	}
}

func (app *appEnv) SinglePrompt() {
	app.loadConfig()
	app.isSinglePrompt = true
	app.currentSession = types.Session{Messages: []types.Message{}}
	app.currentSession.Messages = append(app.currentSession.Messages, createMessage("user", app.InitialPrompt))
	app.singlePrompt()
}

func (app *appEnv) ListAndSelectSession() {
	app.loadConfig()
	app.listSessions = true
	app.sessions = db.GetLastTenSessions()

	sessionNames := []string{}
	for _, e := range app.sessions {
		var name string

		if e.Messages[0].Role == "system" {
			name = e.Messages[1].Content
		} else {
			name = e.Messages[0].Content
		}

		if len(name) > 90 {
			name = strings.TrimSpace(name[:90]) + "..."
		}

		sessionNames = append(sessionNames, name)
	}

	selectSessionPromptContent := promptSelectContent{
		label:        "Select a previous chat",
		selectValues: sessionNames,
	}
	promptResult := promptGetSelect(selectSessionPromptContent)

	app.currentSession = app.sessions[promptResult.index]
	app.sessions = nil
	for _, e := range app.currentSession.Messages {
		if e.Role == "user" {
			fmt.Println("USER: ", e.Content+"\n")
		} else if e.Role == "system" {
			fmt.Println("SYSTEM: ", e.Content+"\n")
		} else {
			printResponse(e.Content + "\n")
		}
	}
}

func (app *appEnv) GenerateImage() {
	fmt.Print(clearScreen)
	req := buildImageRequest(app)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	printResponse(stringifyResponseBody(resp))
	print()
}
