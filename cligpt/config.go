package cligpt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

type Config struct {
	Model       string  `yaml:"model"`
	Token       string  `yaml:"token"`
	Personality string  `yaml:"personality"`
	Temperature float64 `yaml:"temperature"`
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
	fmt.Println("Config file created at: ", f.Name())
}

func saveToConfig(key string, value string) {
	path := getConfigPath()

	config := parseConfig()

	switch key {
	case "model":
		config.Model = value
	case "token":
		config.Token = value
	case "personality":
		config.Personality = value
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
