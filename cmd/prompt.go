/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/eitamonya/cligpt/cligpt"

	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Prompt the model with a single prompt",
	Long:  `This command will prompt the model with a single prompt`,
	Run: func(cmd *cobra.Command, args []string) {
		var prompt string
		for _, arg := range args {
			prompt += arg + " "
		}

		isJson, _ := cmd.Flags().GetBool("json")
		app := cligpt.InitApp()
		app.InitialPrompt = prompt
		app.OutputJSON = isJson
		app.SinglePrompt()
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
	promptCmd.Flags().BoolP("json", "j", false, "Use this flag if you want the response to be output in json")
}
