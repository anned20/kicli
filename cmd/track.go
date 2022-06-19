package cmd

import (
	"errors"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Start tracking time",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		// Fetch projects and activities
		activitiesChannel := make(chan []kimai.Activity, 1)

		go func() {
			value, err := kimaiClient.GetActivities()

			if err != nil {
				logrus.Fatal(err)
			}

			activitiesChannel <- value
		}()

		projectsChannel := make(chan []kimai.Project, 1)

		go func() {
			value, err := kimaiClient.GetProjects()

			if err != nil {
				logrus.Fatal(err)
			}

			projectsChannel <- value
		}()

		// Wait for the activities channel to be populated
		activities := <-activitiesChannel

		activityNames := lo.Map(activities, func(activity kimai.Activity, i int) string {
			return activity.Name
		})

		activityName := ""
		prompt := &survey.Select{
			Message: "Choose an activity:",
			Options: activityNames,
			Default: "Softwaredevelopment",
		}
		err := survey.AskOne(prompt, &activityName)

		if err != nil {
			logrus.Fatal(err)
		}

		activity, found := lo.Find(activities, func(project kimai.Activity) bool {
			return project.Name == activityName
		})

		if !found {
			logrus.Fatal("Activity not found")
		}

		// Wait for the projects channel to be populated
		projects := <-projectsChannel

		projectNames := lo.Map(projects, func(project kimai.Project, i int) string {
			return project.Name
		})

		projectNames = append([]string{"New project"}, projectNames...)

		projectName := ""
		prompt = &survey.Select{
			Message: "Choose a project",
			Options: projectNames,
			Default: "New project",
			Description: func(value string, index int) string {
				if value == "New project" {
					return "Create a new project"
				}

				project, found := lo.Find(projects, func(project kimai.Project) bool {
					return project.Name == value
				})

				if !found {
					return "Project not found"
				}

				return project.ParentTitle
			},
		}
		err = survey.AskOne(prompt, &projectName)

		if err != nil {
			logrus.Fatal(err)
		}

		var project kimai.Project

		if projectName != "New project" {
			project, found = lo.Find(projects, func(project kimai.Project) bool {
				return project.Name == projectName
			})

			if !found {
				logrus.Fatal("Project not found")
			}
		} else {
			_project, err := createProject(kimaiClient)

			if err != nil {
				logrus.Fatal(err)
			}

			project = *_project
		}

		// Start tracking
		timesheet := &kimai.Timesheet{
			ProjectID:  project.ID,
			ActivityID: activity.ID,
		}

		_, err = kimaiClient.CreateTimesheet(timesheet)

		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("Tracking %s on %s\n", color.YellowString(activity.Name), color.YellowString(project.Name))
	},
}

func init() {
	rootCmd.AddCommand(trackCmd)
}

func createProject(kimaiClient *kimai.KimaiClient) (*kimai.Project, error) {
	customers, err := kimaiClient.GetCustomers()

	if err != nil {
		return nil, err
	}

	customerNames := lo.Map(customers, func(customer kimai.Customer, i int) string {
		return customer.Name
	})

	customerName := ""
	customerPrompt := &survey.Select{
		Message: "Choose a customer",
		Options: customerNames,
	}
	err = survey.AskOne(customerPrompt, &customerName)

	if err != nil {
		return nil, err
	}

	customer, found := lo.Find(customers, func(customer kimai.Customer) bool {
		return customer.Name == customerName
	})

	if !found {
		return nil, errors.New("Customer not found")
	}

	project := &kimai.Project{
		CustomerID: customer.ID,
	}

	projectPrompt := &survey.Input{
		Message: "Project name:",
	}

	err = survey.AskOne(projectPrompt, &project.Name, survey.WithValidator(survey.MinLength(2)))

	if err != nil {
		return nil, err
	}

	project, err = kimaiClient.CreateProject(project)

	if err != nil {
		return nil, err
	}

	return project, nil
}
