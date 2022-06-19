package cmd

import (
	"fmt"
	"strconv"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// trackStatusCmd represents the track status command
var trackStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of current timesheet",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		// Start tracking
		timesheet, err := kimaiClient.GetActiveTimesheet()

		if err != nil {
			logrus.Fatal(err)
		}

		if timesheet == nil {
			fmt.Println(color.RedString("No timesheet active"))
			return
		}

		fmt.Printf("Active timesheet %s\n", color.YellowString(strconv.Itoa(timesheet.ID)))
		fmt.Printf("Project: %s\n", color.YellowString(timesheet.Project.Name))
		fmt.Printf("Customer: %s\n", color.YellowString(timesheet.Project.Customer.Name))
		fmt.Printf("Activity: %s\n", color.YellowString(timesheet.Activity.Name))
		fmt.Printf("Started: %s\n", color.YellowString(timesheet.Start.String()))
		fmt.Printf("Duration: %s\n", color.YellowString(timesheet.RealDuration().String()))
	},
}

func init() {
	trackCmd.AddCommand(trackStatusCmd)
}
