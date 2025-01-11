package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "notification-service",
	Short: "Notification service CLI",
}

func Execute() error {
	// cobra.CheckErr(rootCmd.Execute())
	return rootCmd.Execute()
}
