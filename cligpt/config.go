package cligpt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	folderName = ".cligpt"
	configName = "config.yaml"
)

var models = map[string]string{
	"chatgpt": "gpt-3.5-turbo",
	"gpt4":    "gpt-4",
}

type Personality struct {
	Name    string `yaml:"name"`
	Active  bool   `yaml:"active"`
	Context string `yaml:"context"`
}

type Config struct {
	Model         string        `yaml:"model"`
	Token         string        `yaml:"token"`
	Personalities []Personality `yaml:"personalities"`
	Temperature   float64       `yaml:"temperature"`
	MaxTokens     int           `yaml:"max_tokens"`
}

func getConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(homedir, folderName, configName)
}

func createConfig() {
	if _, err := os.Stat(getConfigPath()); err == nil {
		log.Default().Println("Config file found, skipping creation...")
		return
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	pathToFolder := filepath.Join(homedir, folderName)

	if _, err := os.Stat(pathToFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathToFolder, 0775)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Create(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// Set default values
	config := Config{
		Model:         "",
		Personalities: []Personality{},
		Temperature:   1.0,
		MaxTokens:     0,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(f.Name(), data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Config file created at: ", f.Name())
}

func saveToConfig(key string, value string) {
	path := getConfigPath()

	config := parseConfig()

	switch key {
	case "temperature":
		res, err := strconv.ParseFloat(value, 64)
		if err != nil {
		}
		config.Temperature = res
	case "max_tokens":
		res, err := strconv.Atoi(value)
		if err != nil {
		}
		config.MaxTokens = res
	case "model":
		config.Model = value
	case "token":
		config.Token = value
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(path, data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(key+" saved to config file at: ", path)
}

func parseConfig() Config {
	path := getConfigPath()

	if _, err := os.Stat(path); err != nil {
		log.Fatal("Config file not found, please run `cligpt init` first")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func SelectAndSaveModel() {
	selectModelPromptContent := promptSelectContent{
		label:        "Select a model",
		selectValues: []string{"chatgpt", "gpt4"},
	}
	promptResult := promptGetSelect(selectModelPromptContent)
	saveToConfig("model", models[promptResult.value])
}

func GetAndSaveToken() {
	getTokenInputContent := promptInputContent{
		errorMsg: "Please enter a valid token",
		label:    "Enter your OpenAI token:",
		isValidInputString: func(input string) bool {
			return strings.HasPrefix(input, "sk-")
		},
	}

	token := promptGetInput(getTokenInputContent)
	saveToConfig("token", token)
}

func AddPersonality() {
	getPersonalityNameInputContent := promptInputContent{
		errorMsg: "Please enter a valid name",
		label:    "Enter a name for the personality:",
		isValidInputString: func(input string) bool {
			return len(input) > 0
		},
	}

	name := promptGetInput(getPersonalityNameInputContent)

	getPersonalityContextInputContent := promptInputContent{
		errorMsg: "Please enter a valid context",
		label:    "Enter a context for the personality:",
		isValidInputString: func(input string) bool {
			return len(input) > 0
		},
	}

	context := promptGetInput(getPersonalityContextInputContent)

	config := parseConfig()

	config.Personalities = append(config.Personalities, Personality{Name: name, Context: context})

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(getConfigPath(), data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Personality saved to config file at: ", getConfigPath())
}

func SetActivePersonality() {
	config := parseConfig()

	if len(config.Personalities) == 0 {
		log.Fatal("No personalities found, please add one first")
	}

	var personalityNames []string
	for _, personality := range config.Personalities {
		personalityNames = append(personalityNames, personality.Name)
	}

	selectPersonalityPromptContent := promptSelectContent{
		label:        "Select a personality",
		selectValues: personalityNames,
	}
	promptResult := promptGetSelect(selectPersonalityPromptContent)

	var selected string
	for i := range config.Personalities {
		if config.Personalities[i].Name == promptResult.value {
			config.Personalities[i].Active = true
			selected = config.Personalities[i].Context
		} else {
			config.Personalities[i].Active = false
		}
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(getConfigPath(), data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Print("Selected personality: ")
	printResponse(selected)
	fmt.Println()
}

func SetTemperature() {
	getTemperatureInputContent := promptInputContent{
		errorMsg: "Please enter a valid temperature",
		label:    "Enter a temperature:",
		isValidInputString: func(input string) bool {
			res, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return false
			}
			return res >= 0 && res <= 1
		},
	}

	temperature := promptGetInput(getTemperatureInputContent)
	saveToConfig("temperature", temperature)
}

func SetMaxTokens() {
	getMaxTokensInputContent := promptInputContent{
		errorMsg: "Please enter a valid max tokens",
		label:    "Enter a max tokens:",
		isValidInputString: func(input string) bool {
			res, err := strconv.Atoi(input)
			if err != nil {
				return false
			}
			return res >= 1
		},
	}

	maxTokens := promptGetInput(getMaxTokensInputContent)
	saveToConfig("max_tokens", maxTokens)
}
