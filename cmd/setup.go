package cmd

import (
	"fmt"
	"os"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	survey "github.com/AlecAivazis/survey/v2"
	yaml "gopkg.in/yaml.v3"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the application",
	Long: `Before using the application, some variables need to be configured.

This command will help you to configure the application. It will ask you for
the following variables:

	- Base URL
	- Username
	- API token
	- Default activity
`,
	Run: func(cmd *cobra.Command, args []string) {
		var config struct {
			BaseURL         string `yaml:"kimai_base_url"`
			Username        string `yaml:"kimai_username"`
			Token           string `yaml:"kimai_api_token"`
			DefaultActivity string `yaml:"kimai_default_activity"`
		}

		// Check if file kicli.yaml exists, if it does not, create it
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			// Create file
			file, err := os.Create(cfgFile)

			if err != nil {
				logrus.Fatal(err)
			}

			defer file.Close()

			logrus.Debugf("Created config file %s", cfgFile)

			// Write default config to file
			err = yaml.NewEncoder(file).Encode(config)

			if err != nil {
				logrus.Fatal(err)
			}
		} else {
			// Read config from file
			file, err := os.Open(cfgFile)

			if err != nil {
				logrus.Fatal(err)
			}

			defer file.Close()

			logrus.Debug("Config file exists, reading it")

			// Decode config from file
			err = yaml.NewDecoder(file).Decode(&config)

			if err != nil {
				logrus.Fatal(err)
			}
		}

		questions := []*survey.Question{
			{
				Name: "baseURL",
				Prompt: &survey.Input{
					Message: "Base URL:",
					Default: config.BaseURL,
				},
			},
			{
				Name: "username",
				Prompt: &survey.Input{
					Message: "Username:",
					Default: config.Username,
				},
			},
			{
				Name: "token",
				Prompt: &survey.Input{
					Message: "API token:",
					Default: config.Token,
				},
			},
			{
				Name: "defaultActivity",
				Prompt: &survey.Input{
					Message: "Default activity (e.g. \"Softwaredevelopment\"):",
					Default: config.DefaultActivity,
				},
			},
		}

		err := survey.Ask(questions, &config)

		if err != nil {
			logrus.Fatal(err)
		}

		// Open file
		file, err := os.OpenFile(cfgFile, os.O_RDWR, 0644)

		if err != nil {
			logrus.Fatal(err)
		}

		defer file.Close()

		kimaiClient := kimai.NewKimaiClient(
			config.BaseURL,
			config.Username,
			config.Token,
		)

		// Get me
		me, err := kimaiClient.GetMe()

		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("You are logged in as %s (%s)\n", color.YellowString(me.Alias), color.YellowString(me.Username))
		fmt.Printf("Config is written to %s\n", color.YellowString(cfgFile))

		// Write default config to file
		err = yaml.NewEncoder(file).Encode(config)

		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
