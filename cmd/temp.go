/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"cligpt/cligpt"

	"github.com/spf13/cobra"
)

// tempCmd represents the temp command
var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Set temperature",
	Long:  `This command will set the temperature`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.SetTemperature()
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)
}
