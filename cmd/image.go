/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/eitamonya/cligpt/cligpt"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate a DALL-E image using the OpenAI API",
	Long: `Usage:
	cligpt image [prompt]

	Generate a DALL-E image using the OpenAI API.
	Please note that this is charged on different basis compared to the ChatGPT/GPT-4 API.`,
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string
		for _, arg := range args {
			prompt += arg + " "
		}

		app := cligpt.InitApp()
		app.InitialPrompt = prompt
		app.GenerateImage()
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
