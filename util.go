package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Model string `yaml:"model"`
	Token string `yaml:"token"`
}

func getConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(homedir, ".cligpt")
	return path
}

func createConfig() {
	path := getConfigPath()

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	fmt.Println("Config file created at: ", f.Name())
}

func saveToken(token string) {
	path := getConfigPath()

	config := parseConfig()

	config.Token = token

	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(path, data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("Token saved to config file at: ", path)
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func showSpinner(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	fmt.Print(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "Loading"))
	for {
		select {
		case <-ticker.C:
			fmt.Print(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, "."))
		case <-shutdownCh:
			return
		}
	}
}

func getUserInput() string {
	fmt.Print("> ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}
