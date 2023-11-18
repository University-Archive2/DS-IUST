package main

import (
	"Task/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server_start",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
