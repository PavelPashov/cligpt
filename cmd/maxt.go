/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"cligpt/cligpt"

	"github.com/spf13/cobra"
)

// maxtCmd represents the maxt command
var maxtCmd = &cobra.Command{
	Use:   "maxt",
	Short: "Set max token usage",
	Long:  `This command will set the max token usage`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.SetMaxTokens()
	},
}

func init() {
	rootCmd.AddCommand(maxtCmd)
}
