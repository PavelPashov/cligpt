package cligpt

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"cligpt/db"
	"cligpt/types"
)

const clearScreen string = "\033[H\033[2J"
const responseColor string = "\x1b[%dm%s\x1b[0m"

var models = map[string]string{
	"chatgpt": "gpt-3.5-turbo",
	"gpt4":    "gpt-4",
}

func CLI(args []string) {
	var app appEnv
	app.fromArgs(args)
	app.run()
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

type appEnv struct {
	messages       []types.Message
	model          string
	token          string
	outputJSON     bool
	isSinglePrompt bool
	initialPrompt  string
	temperature    float64
	personality    string
	listSessions   bool
	sessions       []types.Session
	currentSession types.Session
}

func (app *appEnv) getDefaultConfig(fl *flag.FlagSet) {
	config := parseConfig()

	if app.token == "" && config.Token == "" {
		fl.Usage()
		log.Fatal("Token not provided nor found in config!")
	} else if app.token == "" {
		app.token = config.Token
	}

	if app.model == "" && config.Model == "" {
		app.model = models["chatgpt"]
	} else {
		app.model = config.Model
	}

	app.personality = config.Personality
	app.temperature = config.Temperature
}

func isFlagPassed(fl *flag.FlagSet, name string) bool {
	found := false
	fl.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func (app *appEnv) fromArgs(args []string) {
	app.messages = []types.Message{}

	fl := flag.NewFlagSet("cli-gpt", flag.ContinueOnError)

	fl.StringVar(
		&app.model, "m", "", "AI model: (chatgpt or gpt4)",
	)

	fl.StringVar(
		&app.token, "t", "", "OpenAI Token (sk-...)",
	)

	fl.BoolVar(
		&app.outputJSON, "j", false, "Use this flag if you want the response to be output in json",
	)

	fl.BoolVar(
		&app.isSinglePrompt, "s", false, "Use this flag if you want to input a single prompt",
	)

	fl.BoolVar(
		&app.listSessions, "l", false, "Use this flag to list your latest 10 sessions",
	)

	fl.Parse(args)

	for _, arg := range fl.Args() {
		app.initialPrompt += arg + " "
	}

	if isFlagPassed(fl, "m") {
		fmt.Println(args)

		selectedModel := models[strings.ToLower(app.model)]

		if selectedModel == "" {
			fl.Usage()
			log.Fatal("Invalid model provided: ", app.model)
		}
		saveToConfig("model", selectedModel)
	}

	if isFlagPassed(fl, "t") {
		if app.token == "" || !strings.HasPrefix(app.token, "sk-") {
			fl.Usage()
			log.Fatal("Please provide a valid token")
		}
		saveToConfig("token", app.token)
	}

	if !app.isSinglePrompt && app.outputJSON {
		fl.Usage()
		log.Fatal("Json output only available for single prompt")
	}

	if app.listSessions && (app.isSinglePrompt || app.outputJSON) {
		fl.Usage()
		log.Fatal("Cannot list session in single prompt mode")
	}

	if app.listSessions {
		app.sessions = db.GetLastTenSessions()
	}

	app.getDefaultConfig(fl)
}

func (app *appEnv) printSessions() {
	index := 0
	printResponse("Your latest 10 sessions:\n")
	for _, session := range app.sessions {
		fmt.Print("ID: ", index)
		fmt.Println("  | ", session.Messages[0].Content)
		index++
	}
	printResponse("Please select a session by ID\n")
}

func (app *appEnv) loadSession() {
	for true {
		app.printSessions()
		selectedSession := getUserInput()
		index, err := strconv.Atoi(selectedSession)

		if err != nil || index > len(app.sessions) || index < 0 {
			fmt.Println("Invalid session provided: ", selectedSession)
		} else {
			app.currentSession = app.sessions[index]

			printResponse("ID: " + selectedSession + " | " + strings.Trim(app.currentSession.Messages[0].Content, " ") + " selected")
			println()

			app.messages = app.currentSession.Messages
			break
		}
	}

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

	if app.outputJSON {
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

func (app *appEnv) run() {
	if app.personality != "" {
		app.messages = append(app.messages, createMessage("system", app.personality))
	}

	if app.isSinglePrompt {
		app.messages = append(app.messages, createMessage("user", app.initialPrompt))
		app.singlePrompt()
	} else {
		if app.listSessions {
			app.loadSession()
		}

		for true {
			var input string
			if app.initialPrompt != "" {
				input = app.initialPrompt
				app.initialPrompt = ""
			} else {
				input = getUserInput()
			}
			app.messages = append(app.messages, createMessage("user", input))
			app.sessionPrompt()
		}
	}
}
