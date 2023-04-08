/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/eitamonya/cligpt/cligpt"

	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Update the token",
	Long:  `This command will update the token`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.GetAndSaveToken()
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
