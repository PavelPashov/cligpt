/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/eitamonya/cligpt/cligpt"

	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var personaCmd = &cobra.Command{
	Use:   "persona",
	Short: "Set active persona",
	Long:  `This command will set the active persona`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.SetActivePersonality()
	},
}

var addPersonaCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new persona",
	Long:  `This command will add a new persona`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.AddPersonality()
	},
}

func init() {
	rootCmd.AddCommand(personaCmd)
	personaCmd.AddCommand(addPersonaCmd)
}
