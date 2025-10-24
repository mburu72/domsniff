package validate

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/fatih/color"
	"github.com/mburu72/domsnif/internal/utils"
)

type Result struct {
	Domain string
	Email  string
}

var VerifyFunc = func(v *emailverifier.Verifier, email string) (*emailverifier.Result, error) {
	return v.Verify(email)
}

func RunVerification(inputPath, outputPath string, workers int) error {
	start := time.Now()

	domains, err := readCSV(inputPath)
	if err != nil {
		return fmt.Errorf("read CSV: %w", err)
	}
	fmt.Printf("\n%s Loaded %d domains from %s\n\n", "", len(domains), inputPath)

	verifier := emailverifier.NewVerifier().EnableDomainSuggest().
		EnableSMTPCheck().
		EnableCatchAllCheck()
	validFile := "data/verified_emails.csv"
	failFile := "data/failed_v_emails.csv"

	var wg sync.WaitGroup
	domainChan := make(chan string)
	resultChan := make(chan []Result)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainChan {
				emails := guessEmails(domain)
				valids := []Result{}
				for _, e := range emails {
					ret, err := VerifyFunc(verifier, e)
					fmt.Printf("Verify result: ret=%v, err=%v\n", ret, err)
					fmt.Println(ret.Email, err)

					if err != nil || ret == nil {
						if ret == nil {
							fmt.Printf("Error: ret is nil for %s\n", e)
						} else {
							fmt.Printf("Error verifying %s: %v\n", e, err)
						}
						continue
					}

					if isDeliverable(ret, err) {
						color.Green("%s", e)
						valids = append(valids, Result{Domain: domain, Email: e})
						utils.AppendToCsv(validFile, [][]string{{domain, e}})
					} else {
						failReason := "unknown"
						if err != nil {
							failReason = err.Error()
						}
						color.Red("%s", e)
						utils.AppendToCsv(failFile, [][]string{{domain, e, failReason}})
					}

				}
				if len(valids) > 0 {
					resultChan <- valids
				}
			}
		}()
	}
	go func() {
		for _, d := range domains {
			domainChan <- d
		}
		close(domainChan)
		wg.Wait()
		close(resultChan)
	}()
	allResults := [][]string{{"domain", "email"}}
	for res := range resultChan {
		for _, r := range res {
			allResults = append(allResults, []string{r.Domain, r.Email})
		}
	}
	fmt.Println()
	if len(allResults) > 1 {
		saveCSV(allResults, outputPath)
		fmt.Printf("\n Saved %d verified emails to %s in %s\n", len(allResults)-1,
			outputPath, time.Since(start).Round(time.Second))
	} else {
		fmt.Printf("\n No valid emails found.\n")
	}
	return nil
}
func isDeliverable(ret *emailverifier.Result, err error) bool {
	if err != nil || ret == nil {
		return false
	}

	if !ret.Syntax.Valid {
		return false
	}

	if !ret.HasMxRecords {
		fmt.Printf("Domain %s has no MX records\n", ret.Syntax.Domain)
		return false
	}

	if ret.SMTP != nil {
		if ret.SMTP.Deliverable {
			return true
		}
		if ret.SMTP.CatchAll {
			fmt.Printf("Domain %s is catch-all — treating as valid\n", ret.Syntax.Domain)
			return true
		}
		fmt.Printf("SMTP check failed for %s: deliverable=%v hostExists=%v fullInbox=%v disabled=%v\n",
			ret.Email,
			ret.SMTP.Deliverable,
			ret.SMTP.HostExists,
			ret.SMTP.FullInbox,
			ret.SMTP.Disabled,
		)
		return false
	}

	fmt.Printf("SMTP check not performed for %s\n", ret.Email)
	return false
}

func readCSV(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var domains []string

	for _, I := range lines {
		if len(I) > 0 && I[0] != "domain" {
			domains = append(domains, strings.TrimSpace(I[0]))
		}
	}
	return domains, nil
}

func saveCSV(records [][]string, path string) {
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create CSV: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	for _, rec := range records {
		_ = w.Write(rec)
	}
}

func guessEmails(domain string) []string {
	prefixes := []string{"info", "contact", "support", "admin", "sales", "team"}
	emails := make([]string, len(prefixes))
	for i, p := range prefixes {
		emails[i] = fmt.Sprintf("%s@%s", p, domain)
	}
	return emails
}
