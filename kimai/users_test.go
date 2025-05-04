package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetMe(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/users/me",
		httpmock.NewStringResponder(200, `
			{
				"language": "en",
				"timezone": "Europe/Berlin",
				"preferences": [],
				"id": 1,
				"username": "admin",
				"enabled": true,
				"roles": [
					"admin"
				],
				"alias": "admin",
				"title": "admin",
				"avatar": "",
				"account_number": "",
				"color": ""
			}
		`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin")

	me, err := client.GetMe()

	assert.NoError(t, err)

	assert.Equal(t, 1, me.ID)
	assert.Equal(t, "admin", me.Username)
}

func TestGetMalformedMe(t *testing.T) {
	cases := map[string]struct {
		response  string
		errString string
		meNil     bool
	}{
		"malformed me": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString: "Malformed",
			meNil:     true,
		},
		"malformed json": {
			response: `[{
				"randomKey": "randomValue"
			}]`,
			errString: "json: cannot unmarshal",
			meNil:     true,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/users/me",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin")

			me, err := client.GetMe()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.meNil {
				assert.Nil(t, me)
			} else {
				assert.NotNil(t, me)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetMeWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin")

	me, err := client.GetMe()

	assert.Error(t, err)
	assert.Nil(t, me)
}
