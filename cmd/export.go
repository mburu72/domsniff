package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var input string
var format string

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&input, "input", "i", "validated.json", "Validated JSON input")
	exportCmd.Flags().StringVar(&format, "format", "csv", "Export format (csv|json)")
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export validation results",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, _ := os.ReadFile(input)
		var results []map[string]string
		json.Unmarshal(data, &results)

		if format == "csv" {
			file, _ := os.Create("results.csv")
			defer file.Close()
			writer := csv.NewWriter(file)
			writer.Write([]string{"domain", "email", "status"})
			for _, r := range results {
				writer.Write([]string{r["domain"], r["email"], r["status"]})
			}
			writer.Flush()
			fmt.Println("Exported to results.csv")
		} else {
			os.WriteFile("results.json", data, 0644)
			fmt.Println("Exported to results.json")
		}
		return nil
	},
}
