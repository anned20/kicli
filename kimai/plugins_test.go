package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPlugins(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/plugins",
		httpmock.NewStringResponder(200, `[
			{
				"name": "Plugin 1",
				"version": "1.0.0"
			},
			{
				"name": "Plugin 2",
				"version": "0.0.1-alpha-2"
			}
		]`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	plugins, err := client.GetPlugins()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(plugins))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a list of plugins is returned
	assert.Equal(t, "Plugin 1", plugins[0].Name)
	assert.Equal(t, "1.0.0", plugins[0].Version)
	assert.Equal(t, "Plugin 2", plugins[1].Name)
	assert.Equal(t, "0.0.1-alpha-2", plugins[1].Version)
}

func TestGetMalformedPlugins(t *testing.T) {
	cases := map[string]struct {
		response   string
		errString  string
		pluginsNil bool
	}{
		"malformed plugins": {
			response: `[
				{
					"randomKey": "randomValue"
				}
			]`,
			errString:  "Malformed",
			pluginsNil: true,
		},
		"malformed json": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:  "json: cannot unmarshal",
			pluginsNil: true,
		},
		"no plugins": {
			response:   `[]`,
			pluginsNil: false,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/plugins",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin", "admin")

			plugins, err := client.GetPlugins()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.pluginsNil {
				assert.Nil(t, plugins)
			} else {
				assert.NotNil(t, plugins)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetPluginsWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin", "admin")

	plugins, err := client.GetPlugins()

	assert.Error(t, err)
	assert.Equal(t, 0, len(plugins))
}
