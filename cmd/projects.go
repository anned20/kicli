package cmd

import (
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Project related commands",
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
