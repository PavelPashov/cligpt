package cligpt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	folderName = ".cligpt"
	configName = "config.yaml"
)

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

func saveToConfig(key string, value string ) {
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

	fmt.Println(key + " saved to config file at: ", path)
}

func saveConfig(newConfig Config) {
	data, err := yaml.Marshal(&newConfig)
	if err != nil {
		log.Fatal(err)
	}

	path := getConfigPath()

	err2 := ioutil.WriteFile(path, data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func parseConfig() Config {
	path := getConfigPath()

	if _, err := os.Stat(path); err != nil {
		createConfig()
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
