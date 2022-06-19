package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cfgFile holds the path to the config file
var cfgFile string

// logLevel is the log level which can be one of "debug", "info", "warn", "error"
var logLevel string

// noColor is a flag to disable color output
var noColor bool

type contextKey int

const (
	KimaiClientKey contextKey = iota
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kicli",
	Short: "CLI to control a Kimai server",
	Long: `This CLI is used to be a companion to a Kimai server.

It allows you to manage your projects, tasks and timesheets.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd != setupCmd {
			required := []string{"kimai_base_url", "kimai_username", "kimai_api_token"}

			missing := []string{}

			for _, key := range required {
				if viper.GetString(key) == "" {
					missing = append(missing, key)
				}
			}

			if len(missing) > 0 {
				fmt.Println(color.New(color.FgRed).Add(color.Bold).Sprintf("Missing configuration variables: %s", strings.Join(missing, ", ")))
				fmt.Printf("Please run %s to configure the CLI.\n", color.New(color.FgYellow).Sprint("kicli setup"))
				os.Exit(1)
			}

			ctx := context.Background()

			// Initialize the Kimai client
			kimaiClient := kimai.NewKimaiClient(
				viper.GetString("kimai_base_url"),
				viper.GetString("kimai_username"),
				viper.GetString("kimai_api_token"),
			)

			// Add the kimai client to the context
			ctx = context.WithValue(ctx, KimaiClientKey, kimaiClient)

			// Set the context to the command
			cmd.SetContext(ctx)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kicli.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set log level
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kicli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kicli")

		cfgFile = fmt.Sprintf("%s/.kicli.yaml", home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set default values for configuration variables
	viper.SetDefault("kimai_base_url", "")
	viper.SetDefault("kimai_username", "")
	viper.SetDefault("kimai_api_token", "")

	// Set the color output
	if noColor {
		color.NoColor = true
	}
}
