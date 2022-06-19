package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/version",
		httpmock.NewStringResponder(200, `{
			"version": "1.0.0"
		}`),
	)

	clientWithoutTrailingSlash := NewKimaiClient("http://localhost:8080", "admin", "admin")
	clientWithTrailingSlash := NewKimaiClient("http://localhost:8080/", "admin", "admin")

	version, err := clientWithoutTrailingSlash.GetVersion()
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", version.Version)

	version, err = clientWithTrailingSlash.GetVersion()
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", version.Version)
}

func TestApiError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/version",
		httpmock.NewStringResponder(500, `{
			"error": "Internal Server Error"
		}`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	_, err := client.GetVersion()

	assert.Error(t, err)
	assert.ErrorContains(t, err, "Internal Server Error")
}
