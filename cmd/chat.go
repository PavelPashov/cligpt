/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/eitamonya/cligpt/cligpt"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat with the model",
	Long:  `This command will start a chat with the model, you can specify the initial prompt to use with the --prompt flag`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		if prompt == "" && len(args) > 0 {
			for _, arg := range args {
				prompt += arg + " "
			}
		}

		app := cligpt.InitApp()
		app.InitialPrompt = prompt
		app.Chat()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the saved chat sessions",
	Long:  `This command will list the saved chat sessions`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		app := cligpt.InitApp()
		app.InitialPrompt = prompt
		app.ListAndSelectSession()
		app.Chat()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringP("prompt", "p", "", "The initial prompt to use for the chat session\nUsage: --prompt \"Hello, how are you?\"")
	chatCmd.AddCommand(listCmd)
}
