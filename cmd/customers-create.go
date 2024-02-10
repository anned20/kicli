package cmd

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/anned20/kicli/kimai"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// customersCreateCmd represents the track command
var customersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new customer",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		kimaiClient := ctx.Value(KimaiClientKey).(*kimai.KimaiClient)

		customer := &kimai.Customer{}

		customerNamePrompt := &survey.Input{
			Message: "Customer name:",
		}

		err := survey.AskOne(customerNamePrompt, &customer.Name, survey.WithValidator(survey.MinLength(2)))

		if err != nil {
			logrus.Fatal(err)
		}

		defaultCustomerCountry := viper.GetString("kimai_default_country")

		customerCountryPrompt := &survey.Input{
			Message: "Customer country:",
			Default: defaultCustomerCountry,
		}

		err = survey.AskOne(customerCountryPrompt, &customer.Country, survey.WithValidator(survey.MinLength(2)))

		if err != nil {
			logrus.Fatal(err)
		}

		defaultCustomerCurrency := viper.GetString("kimai_default_currency")

		customerCurrencyPrompt := &survey.Input{
			Message: "Customer currency:",
			Default: defaultCustomerCurrency,
		}

		err = survey.AskOne(customerCurrencyPrompt, &customer.Currency, survey.WithValidator(survey.MinLength(2)))

		if err != nil {
			logrus.Fatal(err)
		}

		defaultCustomerTimezone := viper.GetString("kimai_default_timezone")

		customerTimezonePrompt := &survey.Input{
			Message: "Customer timezone:",
			Default: defaultCustomerTimezone,
		}

		err = survey.AskOne(customerTimezonePrompt, &customer.Timezone, survey.WithValidator(survey.MinLength(2)))

		if err != nil {
			logrus.Fatal(err)
		}

		customer, err = kimaiClient.CreateCustomer(customer)

		if err != nil {
			logrus.Fatal(err)
		}

		customerIdString := fmt.Sprintf("%d", customer.ID)

		fmt.Printf("Created customer \"%s\" with ID %s", color.YellowString(customer.Name), color.YellowString(customerIdString))
	},
}

func init() {
	customersCmd.AddCommand(customersCreateCmd)
}
