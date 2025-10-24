package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"github.com/mburu72/domsnif/internal/unzip"
)

var unzipandRead = unzip.UnzipAndRead
var httpGet = http.Get

func FetchNRD() ([]string, string, error) {

	c := colly.NewCollector()

	type FileInfo struct {
		Date string
		Link string
	}
	var files []FileInfo

	c.OnHTML("tbody tr", func(h *colly.HTMLElement) {
		date := h.ChildText("td:nth-of-type(3)")
		link := h.ChildAttr("a[href*='/nrd']", "href")

		if link != "" && date != "" {
			files = append(files, FileInfo{
				Date: date,
				Link: link,
			})
		}
	})
	color.Green("Scouting for newly registered domains...")
	c.Visit("https://www.whoisds.com/newly-registered-domains")

	if len(files) == 0 {
		color.Red("No download links found.")
		return nil, "", fmt.Errorf("no download links")
	}

	file := files[0]
	color.Green("Found file with latest domains registered.")
	color.Green("Downloading file....")
	resp, err := httpGet(file.Link)
	if err != nil {
		color.Red("Download failed:", err)
		return nil, "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("non-OK response: %s", resp.Status)
	}

	filename := fmt.Sprintf("newly-registered-domains-%s.zip", file.Date)

	outpath := filepath.Join(".", filename)
	out, err := os.Create(outpath)
	if err != nil {
		return nil, "", fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error saving file: %w", err)
	}
	color.Blue("Downloaded: %s", outpath)
	color.Green("Unzipping & extracting domains...")
	domains, err := unzipandRead(outpath)
	if err != nil {
		color.Red("An error occurred when trying to unzip the file: %s", err)
	}
	color.Green("Domains extracted successfully!")
	return domains, outpath, nil
}
