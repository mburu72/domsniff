package unzip

import (
	"archive/zip"
	"bufio"
	"fmt"
	"strings"
)

func UnzipAndRead(zipPath string) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	var domains []string

	for _, f := range r.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".txt") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open text: %w", err)
		}
		defer rc.Close()

		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				domains = append(domains, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("read txt: %w", err)
		}
	}
	return domains, nil
}
