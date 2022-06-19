package cmd

import (
	"fmt"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about the Kimai server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		version, err := kimaiClient.GetVersion()

		if err != nil {
			logrus.Fatalf("Failed to get version: %s", err)
		}

		plugins, err := kimaiClient.GetPlugins()

		if err != nil {
			logrus.Fatalf("Failed to get plugins: %s", err)
		}

		fmt.Printf("Kimai version: %s\n", color.YellowString(version.Version))

		if len(plugins) > 0 {
			for _, plugin := range plugins {
				fmt.Printf("Plugin: %s, version: %s", color.YellowString(plugin.Name), color.YellowString(plugin.Version))
			}
		} else {
			fmt.Println(color.RedString("No plugins installed"))
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
