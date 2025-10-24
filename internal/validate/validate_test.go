package validate

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuessEmails(t *testing.T) {
	emails := guessEmails("example.com")

	assert.Equal(t, 6, len(emails))
	assert.Contains(t, emails, "info@example.com")
	assert.Contains(t, emails, "sales@example.com")
}

func TestIsDeliverable(t *testing.T) {

	ret := &emailverifier.Result{
		Email: "test@example.com",
		Syntax: emailverifier.Syntax{
			Valid:  true,
			Domain: "example.com",
		},
		HasMxRecords: true,
		SMTP: &emailverifier.SMTP{
			Deliverable: true,
		},
	}
	assert.True(t, isDeliverable(ret, nil))

	ret.HasMxRecords = false
	assert.False(t, isDeliverable(ret, nil))

	ret.HasMxRecords = true
	ret.SMTP = &emailverifier.SMTP{CatchAll: true}
	assert.True(t, isDeliverable(ret, nil))

	ret.Syntax.Valid = false
	assert.False(t, isDeliverable(ret, nil))
}

func TestReadAndSaveCSV(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "domains.csv")

	f, err := os.Create(path)
	require.NoError(t, err)
	w := csv.NewWriter(f)
	_ = w.Write([]string{"domain"})
	_ = w.Write([]string{"example.com"})
	w.Flush()
	f.Close()

	domains, err := readCSV(path)
	require.NoError(t, err)
	assert.Equal(t, []string{"example.com"}, domains)

	output := filepath.Join(tmpDir, "output.csv")
	saveCSV([][]string{{"domain", "email"}, {"example.com", "info@example.com"}}, output)

	data, err := os.ReadFile(output)
	require.NoError(t, err)
	assert.Contains(t, string(data), "info@example.com")
}

func TestRunVerification_Basic(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "domains.csv")
	outputFile := filepath.Join(tmpDir, "verified.csv")

	f, err := os.Create(inputFile)
	require.NoError(t, err)
	_, _ = f.WriteString("domain\nexample.com\n")
	f.Close()

	VerifyFunc = func(v *emailverifier.Verifier, email string) (*emailverifier.Result, error) {
		return &emailverifier.Result{
			Email: email,
			Syntax: emailverifier.Syntax{
				Valid:  true,
				Domain: "example.com",
			},
			HasMxRecords: true,
			SMTP: &emailverifier.SMTP{
				Deliverable: true,
			},
		}, nil
	}
	defer func() { VerifyFunc = nil }()

	err = RunVerification(inputFile, outputFile, 1)
	require.NoError(t, err)

	data, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "info@example.com")
	assert.Contains(t, string(data), "contact@example.com")
}
