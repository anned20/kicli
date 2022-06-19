package cmd

import (
	"fmt"
	"strconv"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// trackStopCmd represents the track status command
var trackStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop active timesheet",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		// Start tracking
		activeTimesheet, err := kimaiClient.GetActiveTimesheet()

		if err != nil {
			logrus.Fatal(err)
		}

		if activeTimesheet == nil {
			fmt.Println(color.RedString("No timesheet active"))
			return
		}

		timesheet, err := kimaiClient.StopTimesheet(activeTimesheet.ID)

		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("Timesheet stopped, ID: %s\n", color.YellowString(strconv.Itoa(timesheet.ID)))
		fmt.Printf("Real duration: %s\n", color.YellowString(timesheet.RealDuration().String()))
		fmt.Printf("Billed duration: %s\n", color.YellowString(timesheet.BilledDuration().String()))

		if timesheet.FixedRate != 0 {
			fmt.Printf("Amount: %s (fixed)", color.New(color.FgYellow).Sprintf("%02.2f", timesheet.FixedRate))
		} else {
			fmt.Printf("Amount: %s", color.New(color.FgYellow).Sprintf("%02.2f", timesheet.Rate))
		}
	},
}

func init() {
	trackCmd.AddCommand(trackStopCmd)
}
