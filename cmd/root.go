package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "echo-server",
	Short: "Echo Server is a simple RESTful API server built with Echo framework",
	Long:  `Echo Server is a simple RESTful API server built with Echo framework in Go. It provides endpoints to manage students.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Echo Server CLI")
		cmd.Help()
	},
}

func Execute() error {

	return rootCmd.Execute()
}
