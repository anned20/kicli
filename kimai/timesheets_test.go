package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetTimesheet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/timesheets/1",
		httpmock.NewStringResponder(200, `{
			"id": 1,
			"project": 1,
			"activity": 1,
			"description": "Timesheet 1",
			"fixedRate": 0,
			"hourlyRate": 0,
			"user": 1,
			"exported": false,
			"billable": false,
			"tags": [],
			"duration": 0,
			"rate": 0,
			"internalRate": 0,
			"begin": "2006-01-02T15:04:05-0700",
			"end": "2006-01-02T15:04:05-0700"
		}`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin")

	timesheet, err := client.GetTimesheet(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a timesheet is returned
	assert.Equal(t, "Timesheet 1", timesheet.Description)
}

func TestGetMalformedTimesheet(t *testing.T) {
	cases := map[string]struct {
		response     string
		errString    string
		timesheetNil bool
	}{
		"malformed timesheet": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:    "Malformed",
			timesheetNil: true,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/timesheets/1",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin")

			timesheet, err := client.GetTimesheet(1)

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.timesheetNil {
				assert.Nil(t, timesheet)
			} else {
				assert.NotNil(t, timesheet)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestCreateTimesheet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"POST",
		"http://localhost:8080/timesheets",
		httpmock.NewStringResponder(200, `{
			"id": 1,
			"project": 1,
			"activity": 1,
			"description": "Timesheet 1",
			"fixedRate": 0,
			"hourlyRate": 0,
			"user": 1,
			"exported": false,
			"billable": false,
			"tags": [],
			"duration": 0,
			"rate": 0,
			"internalRate": 0,
			"begin": "2006-01-02T15:04:05-0700",
			"end": "2006-01-02T15:04:05-0700"
		}`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin")

	timesheet := &Timesheet{
		Description: "Timesheet 1",
		ProjectID:   1,
		ActivityID:  1,
	}

	createdTimesheet, err := client.CreateTimesheet(timesheet)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a timesheet is returned
	assert.Equal(t, "Timesheet 1", createdTimesheet.Description)
	assert.Equal(t, 1, createdTimesheet.ID)
	assert.Equal(t, 1, createdTimesheet.ProjectID)
	assert.Equal(t, 1, createdTimesheet.ActivityID)
}
