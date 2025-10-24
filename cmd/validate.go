package cmd

import (
	"github.com/mburu72/domsnif/internal/validate"
	"github.com/spf13/cobra"
)

var inputFile string
var outputFile string
var workers int
var runVerification = validate.RunVerification

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVarP(&inputFile, "input", "i", "domains.json", "Input JSON file of domains")
	validateCmd.Flags().StringVarP(&outputFile, "output", "o", "data/verified_emails.csv", "output CSV file")
	validateCmd.Flags().IntVar(&workers, "w", 5, "Concurrent workers")
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate possible emails for each domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runVerification(inputFile, outputFile, workers); err != nil {
			return err
		}
		return nil
	},
}
