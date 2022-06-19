package cmd

import (
	"fmt"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// meCmd represents the me command
var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Get information about the current user",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		me, err := kimaiClient.GetMe()

		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("Name: %s\n", color.YellowString(me.Alias))
		fmt.Printf("Username: %s\n", color.YellowString(me.Username))
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
