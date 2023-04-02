package cmd

import (
	"cligpt/cligpt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiate the setup for cli-gpt",
	Long:  `This command will initiate the setup for cli-gpt`,
	Run: func(cmd *cobra.Command, args []string) {
		cligpt.Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
