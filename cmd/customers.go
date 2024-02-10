package cmd

import (
	"github.com/spf13/cobra"
)

// customersCmd represents the customers command
var customersCmd = &cobra.Command{
	Use:   "customers",
	Short: "Customer related commands",
}

func init() {
	rootCmd.AddCommand(customersCmd)
}
