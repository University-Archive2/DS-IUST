package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "distributed IR service",
	Short: "Serves Distributed Information Retrieval System",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()
}
