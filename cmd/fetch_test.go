package cmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(strings.Builder)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	_, err = root.ExecuteC()
	return buf.String(), err
}

func TestFetchCmd_Success_NoFilter(t *testing.T) {

	fetchNRDFunc = func() ([]string, string, error) {
		return []string{"example.com", "test.net",
			"example1.com", "test1.net",
			"example2.com", "test2.net",
			"example3.com", "test3.net",
			"example4.com", "test4.net"}, "dummy.zip", nil
	}

	output, err := executeCommand(rootCmd, "fetch")
	fmt.Println(output)
	require.NoError(t, err)

}

func TestFetchCmd_Success_WithFilter(t *testing.T) {
	fetchNRDFunc = func() ([]string, string, error) {
		return []string{"example.com", "test.net"}, "dummy.zip", nil
	}
	filterDomainsFunc = func(domains []string, checkOnline, checkDev bool) []string {
		return []string{"example.com"}
	}
	saveCsvFunc = func(domains []string, filename string) (string, error) {
		return fmt.Sprintf("/tmp/%s", filename), nil
	}

	output, err := executeCommand(rootCmd, "fetch", "--filter")
	require.NoError(t, err)
	fmt.Println(output)
}

func TestFetchCmd_FetchError(t *testing.T) {
	fetchNRDFunc = func() ([]string, string, error) {
		return nil, "", errors.New("network fail")
	}

	defer func() { recover() }()

	_, err := executeCommand(rootCmd, "fetch")
	assert.Error(t, err)
}
