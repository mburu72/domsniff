package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/mburu72/domsnif/internal/fetch"
	"github.com/mburu72/domsnif/internal/filter"
)

var (
	fetchNRDFunc      = fetch.FetchNRD
	filterDomainsFunc = filter.FilterDomains
	saveCsvFunc       = filter.SaveDomainCsv
	filterFlag        bool
	checkOnline       bool
	checkDev          bool
)

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().BoolVarP(&filterFlag, "filter", "f", false, "enable filtering(TLD, online, under development)")
	fetchCmd.Flags().BoolVarP(&checkOnline, "check-online", "o", false, "filter domains based on online status")
	fetchCmd.Flags().BoolVarP(&checkDev, "check-dev", "d", false, "filter domains that are under development")

}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch newly registered domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		domains, zipfile, err := fetchNRDFunc()

		if err != nil {
			return fmt.Errorf("fetch failed: %s", err)
		}
		color.Green("Got %d domains from %s", len(domains), zipfile)
		sampleCount := 10
		if len(domains) < sampleCount {
			sampleCount = len(domains)
		}
		color.White("Sample: %s", domains[:sampleCount])

		if filterFlag {
			color.Blue("Running filtering....")
			filtered := filterDomainsFunc(domains, checkOnline, checkDev)
			filePath, err := filter.SaveDomainCsv(filtered, "filtered_domains.csv")
			if err != nil {
				color.Red("Save failed: %s", err)

			}
			if checkDev {
				color.Blue("Saving under-development domains to data/under_dev.csv...")
				underDevFilePath, err := saveCsvFunc(filtered, "under_dev.csv")
				if err != nil {
					color.Red("Save to under_dev.csv failed: %s", err)
				} else {
					color.Blue("Done! %d under-development domains saved to %s", len(filtered), underDevFilePath)
				}
			}
			color.Blue("Done! %d filtered domains saved to %s in %s\n", len(filtered), filePath,
				time.Since(start).Round(time.Second))
		} else {
			color.Cyan("Running search without filter")
			color.Magenta("Tip: use --filter to enable filtering")
		}
		return nil
	},
}
