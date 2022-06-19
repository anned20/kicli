package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/version",
		httpmock.NewStringResponder(200, `
			{
				"version": "1.0.0",
				"versionId": 1,
				"candidate": "",
				"semver": "1.0.0",
				"name": "Kimai",
				"copyright": "Copyright (c) 2018-2020 Kimai-Development-Team"
			}
		`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	version, err := client.GetVersion()

	assert.NoError(t, err)

	assert.Equal(t, "1.0.0", version.Version)
	assert.Equal(t, 1, version.VersionId)
	assert.Equal(t, "", version.Candidate)
	assert.Equal(t, "1.0.0", version.Semver)
	assert.Equal(t, "Kimai", version.Name)
	assert.Equal(t, "Copyright (c) 2018-2020 Kimai-Development-Team", version.Copyright)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestGetMalformedVersion(t *testing.T) {
	cases := map[string]struct {
		response   string
		errString  string
		versionNil bool
	}{
		"malformed version": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:  "Malformed",
			versionNil: true,
		},
		"malformed json": {
			response: `[{
				"randomKey": "randomValue"
			}]`,
			errString:  "json: cannot unmarshal",
			versionNil: true,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/version",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin", "admin")

			version, err := client.GetVersion()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.versionNil {
				assert.Nil(t, version)
			} else {
				assert.NotNil(t, version)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetVersionWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin", "admin")

	version, err := client.GetVersion()

	assert.Error(t, err)
	assert.Nil(t, version)
}
