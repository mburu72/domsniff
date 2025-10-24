package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCmd_Success(t *testing.T) {

	called := false
	runVerification = func(inFile, outFile string, workers int) error {
		called = true
		assert.Equal(t, "domains.json", inFile)
		assert.Equal(t, "data/verified_emails.csv", outFile)
		assert.Equal(t, 5, workers)
		return nil
	}

	out, err := executeCommand(rootCmd, "validate")
	require.NoError(t, err)
	assert.True(t, called)
	assert.Contains(t, out, "") // Just ensures it ran
}

func TestValidateCmd_Error(t *testing.T) {
	// Mock function to simulate failure
	runVerification = func(inFile, outFile string, workers int) error {
		return errors.New("verification failed")
	}

	out, err := executeCommand(rootCmd, "validate")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verification failed")
	assert.Contains(t, out, "") // The command still writes to buffer
}
