package filter

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

func FilterDomains(domains []string, checkOnline bool, checkDev bool) []string {

	tlds := []string{".com", ".io", ".net", ".co"}
	var filtered []string
	color.Yellow("Filtering domains to: %s", strings.Join(tlds, " "))
	for _, d := range domains {
		for _, t := range tlds {
			if strings.HasSuffix(strings.ToLower(d), t) {

				if checkOnline && !isOnline(d) {
					continue
				}
				if checkDev && !isInDev(d) {
					continue
				}
				filtered = append(filtered, d)
				break
			}
		}
	}
	return filtered
}

func isOnline(domain string) bool {

	client := http.Client{Timeout: 3 * time.Second}
	url := "http://" + domain
	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400
}

func isInDev(domain string) bool {

	client := http.Client{Timeout: 6 * time.Second}
	url := "http://" + domain

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return false
	}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)

	keywords := []string{
		"coming soon", "under construction", "launching soon", "maintenance mode",
		"site under development", "website is being built",
	}
	readBytes := 0
	maxBytes := 100000

	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		readBytes += len(line)
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				color.Green("%s appears to be under development", domain)
				return true
			}
		}
		if readBytes > maxBytes {
			break
		}
	}
	return false
}

func SaveDomainCsv(domains []string, filename string) (string, error) {

	if len(domains) == 0 {
		return "", fmt.Errorf("no domains to save")
	}
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	if filename == "" {
		filename = "domains.csv"
	}

	filepath := filepath.Join("data", filename)
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, d := range domains {
		if err := writer.Write([]string{d}); err != nil {
			return "", fmt.Errorf("write row: %w", err)
		}
	}
	fmt.Printf("Saved %d domains to %s\n", len(domains), filepath)
	return filepath, nil
}
