package cligpt

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
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
	messages       []Message
	model          string
	token          string
	outputJSON     bool
	isSinglePrompt bool
	initialPrompt  string
}

func (app *appEnv) getDefaultConfig(fl *flag.FlagSet) {
	if app.token == "" || app.model == "" {
		config := parseConfig()
		if app.token == "" && config.Token == "" {
			fl.Usage()
			log.Fatal("Token not provided nor found in config!")
		} else {
			app.token = config.Token
		}
		if app.model == "" && config.Model == "" {
			app.model = models["chatgpt"]
		} else {
			app.model = config.Model
		}
	}
}

func (app *appEnv) fromArgs(args []string) {
	// This can later be used for system messages
	app.messages = []Message{}

	fl := flag.NewFlagSet("cli-gpt", flag.ContinueOnError)

	fl.StringVar(
		&app.model, "m", "", "AI model: (chatgpt or gpt4)",
	)

	fl.StringVar(
		&app.model, "t", "", "OpenAI Token (sk-...)",
	)

	fl.BoolVar(
		&app.outputJSON, "j", false, "Use this flag if you want the response to be output in json",
	)

	fl.BoolVar(
		&app.isSinglePrompt, "s", false, "Use this flag if you want to input a single prompt",
	)

	fl.Parse(args)

	for _, arg := range fl.Args() {
		app.initialPrompt += arg + " "
	}

	if isFlagPassed("m") {
		if models[strings.ToLower(app.model)] == "" {
			fl.Usage()
			log.Fatal("Invalid model provided: ", app.model)
		}
	}

	if isFlagPassed("t") {
		if app.token == "" || !strings.HasPrefix(app.token, "sk-") {
			fl.Usage()
			log.Fatal("Please provide a valid token")
		}
		saveToken(app.token)
	}

	if app.isSinglePrompt && isFlagPassed("j") {
		fl.Usage()
		log.Fatal("Json output only available for single prompt")
	}

	if app.isSinglePrompt && app.initialPrompt == "" {
		fl.Usage()
		log.Fatal("Single prompt requires an input - cligpt -s \"<prompt>\"")
	}

	app.getDefaultConfig(fl)
}

func getUserInput() string {
	fmt.Print("> ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func createMessage(role string, content string) Message {
	return Message{Role: role, Content: content}
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
		return
	}

	responseBody := parseResponse(resp)

	printResponse(responseBody.Choices[0].Message.Content)
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

	app.messages = append(app.messages, Message{Role: "assistant", Content: content})
	fmt.Println()
}

func (app *appEnv) run() {
	if app.isSinglePrompt {
		app.messages = append(app.messages, createMessage("user", app.initialPrompt))
		app.singlePrompt()
	} else {
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
