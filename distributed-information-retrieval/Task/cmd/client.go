package main

import (
	"Task/internal/client"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client_start",
	Short: "Starts a client",
	Run: func(cmd *cobra.Command, args []string) {
		client.StartClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
