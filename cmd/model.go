/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"cligpt/cligpt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Change the model configuration",
	Long:  `This command will change the model configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.SelectAndSaveModel()
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
