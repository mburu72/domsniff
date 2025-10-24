package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

var mockUnzipFunc func(string) ([]string, error)

func init() {

	unzipandRead = func(path string) ([]string, error) {
		return mockUnzipFunc(path)
	}
}

func createTempFile(t *testing.T, data []byte) string {
	tmpfile := filepath.Join(t.TempDir(), "test.zip")
	err := os.WriteFile(tmpfile, data, 0644)
	require.NoError(t, err)
	return tmpfile
}

func TestFetchNRD_Success(t *testing.T) {

	respBody := io.NopCloser(bytes.NewBuffer([]byte("fake zip data")))
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       respBody,
		}, nil
	}

	mockUnzipFunc = func(path string) ([]string, error) {
		return []string{"example.com", "test.net"}, nil
	}

	domains, outpath, err := FetchNRD()

	require.NoError(t, err)
	assert.NotEmpty(t, outpath)
	assert.Contains(t, domains, "example.com")
	assert.Contains(t, domains, "test.net")
}

func TestFetchNRD_NoLinksFound(t *testing.T) {
	httpGet = func(url string) (*http.Response, error) {
		return nil, fmt.Errorf("should not be called")
	}

	unzipandRead = func(path string) ([]string, error) {
		return nil, fmt.Errorf("should not be called")
	}

	domains, outpath, err := FetchNRD()

	assert.Error(t, err)
	assert.Nil(t, domains)
	assert.Empty(t, outpath)
}

func TestFetchNRD_DownloadError(t *testing.T) {

	httpGet = func(url string) (*http.Response, error) {
		return nil, fmt.Errorf("network down")
	}

	domains, outpath, err := FetchNRD()

	assert.Error(t, err)
	assert.Nil(t, domains)
	assert.Empty(t, outpath)
}
