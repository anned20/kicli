package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetActivities(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/activities",
		httpmock.NewStringResponder(200, `[
			{
				"id": 1,
				"name": "Activity 1",
				"project": 1,
				"parentTitle": "Project 1",
				"comment": "",
				"visible": true,
				"metaFields": [],
				"budget": 0,
				"timeBudget": 0,
				"budgetType": "",
				"color": ""
			}, {
				"id": 2,
				"name": "Activity 2",
				"project": 1,
				"parentTitle": "Project 1",
				"comment": "",
				"visible": true,
				"metaFields": [],
				"budget": 0,
				"timeBudget": 0,
				"budgetType": "",
				"color": ""
			}
		]`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin")

	activities, err := client.GetActivities()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(activities))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a list of activities is returned
	assert.Equal(t, "Activity 1", activities[0].Name)
	assert.Equal(t, 1, activities[0].ProjectID)
	assert.Equal(t, "Project 1", activities[0].ParentTitle)

	assert.Equal(t, "Activity 2", activities[1].Name)
	assert.Equal(t, 1, activities[1].ProjectID)
	assert.Equal(t, "Project 1", activities[1].ParentTitle)
}

func TestGetMalformedActivities(t *testing.T) {
	cases := map[string]struct {
		response      string
		errString     string
		activitiesNil bool
	}{
		"malformed activities": {
			response: `[
				{
					"randomKey": "randomValue"
				}
			]`,
			errString:     "Malformed",
			activitiesNil: true,
		},
		"malformed json": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:     "json: cannot unmarshal",
			activitiesNil: true,
		},
		"no activities": {
			response:      `[]`,
			activitiesNil: false,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/activities",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin")

			activities, err := client.GetActivities()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.activitiesNil {
				assert.Nil(t, activities)
			} else {
				assert.NotNil(t, activities)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetActivitiesWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin")

	activities, err := client.GetActivities()

	assert.Error(t, err)
	assert.Equal(t, 0, len(activities))
}
