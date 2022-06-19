package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCustomers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/customers",
		httpmock.NewStringResponder(200, `[
			{
				"id": 1,
				"name": "Customer 1",
				"comment": "",
				"visible": true,
				"metaFields": [],
				"color": ""
			}, {
				"id": 2,
				"name": "Customer 2",
				"comment": "Comment 2",
				"visible": false,
				"metaFields": [],
				"color": ""
			}
		]`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	customers, err := client.GetCustomers()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(customers))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a list of customers is returned
	assert.Equal(t, "Customer 1", customers[0].Name)
	assert.Equal(t, 1, customers[0].ID)
	assert.Equal(t, "", customers[0].Comment)
	assert.Equal(t, true, customers[0].Visible)

	assert.Equal(t, "Customer 2", customers[1].Name)
	assert.Equal(t, 2, customers[1].ID)
	assert.Equal(t, "Comment 2", customers[1].Comment)
	assert.Equal(t, false, customers[1].Visible)
}

func TestGetMalformedCustomers(t *testing.T) {
	cases := map[string]struct {
		response     string
		errString    string
		customersNil bool
	}{
		"malformed customers": {
			response: `[
				{
					"randomKey": "randomValue"
				}
			]`,
			errString:    "Malformed",
			customersNil: true,
		},
		"malformed json": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:    "json: cannot unmarshal",
			customersNil: true,
		},
		"no customers": {
			response:     `[]`,
			customersNil: false,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/customers",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin", "admin")

			customers, err := client.GetCustomers()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.customersNil {
				assert.Nil(t, customers)
			} else {
				assert.NotNil(t, customers)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetCustomersWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin", "admin")

	customers, err := client.GetCustomers()

	assert.Error(t, err)
	assert.Equal(t, 0, len(customers))
}
