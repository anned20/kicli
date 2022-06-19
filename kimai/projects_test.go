package kimai

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

/*

	ID          int         `json:"id"`
	Name        string      `json:"name"`
	CustomerID  int         `json:"customer"`
	ParentTitle string      `json:"parentTitle"`
	Comment     string      `json:"comment"`
	Visible     bool        `json:"visible"`
	MetaFields  []MetaField `json:"metaFields"`

	RawStart string `json:"start"`
	RawEnd   string `json:"end"`
*/

func TestGetProjects(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/projects",
		httpmock.NewStringResponder(200, `[
			{
				"id": 1,
				"name": "Project 1",
				"customer": 1,
				"parentTitle": "",
				"comment": "",
				"visible": true,
				"metaFields": [],
				"start": "2006-01-02T15:04:05-0700",
				"end": "2006-01-02T15:04:05-0700"
			},
			{
				"id": 2,
				"name": "Project 2",
				"customer": 1,
				"parentTitle": "",
				"comment": "",
				"visible": true,
				"metaFields": [],
				"start": "2006-01-02T15:04:05-0700",
				"end": "2006-01-02T15:04:05-0700"
			}
		]`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	projects, err := client.GetProjects()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(projects))
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a list of projects is returned
	assert.Equal(t, "Project 1", projects[0].Name)
	assert.Equal(t, "Project 2", projects[1].Name)
}

func TestGetMalformedProjects(t *testing.T) {
	cases := map[string]struct {
		response    string
		errString   string
		projectsNil bool
	}{
		"malformed projects": {
			response: `[
				{
					"randomKey": "randomValue"
				}
			]`,
			errString:   "Malformed",
			projectsNil: true,
		},
		"malformed json": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:   "json: cannot unmarshal",
			projectsNil: true,
		},
		"no projects": {
			response:    `[]`,
			projectsNil: false,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/projects",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin", "admin")

			projects, err := client.GetProjects()

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.projectsNil {
				assert.Nil(t, projects)
			} else {
				assert.NotNil(t, projects)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestGetProjectsWithNoAPI(t *testing.T) {
	client := NewKimaiClient("http://localhost:65500", "admin", "admin")

	projects, err := client.GetProjects()

	assert.Error(t, err)
	assert.Equal(t, 0, len(projects))
}

func TestGetProject(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8080/projects/1",
		httpmock.NewStringResponder(200, `{
			"id": 1,
			"name": "Project 1",
			"customer": 1,
			"parentTitle": "",
			"comment": "",
			"visible": true,
			"metaFields": [],
			"start": "2006-01-02T15:04:05-0700",
			"end": "2006-01-02T15:04:05-0700"
		}`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	project, err := client.GetProject(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a project is returned
	assert.Equal(t, "Project 1", project.Name)
}

func TestGetMalformedProject(t *testing.T) {
	cases := map[string]struct {
		response   string
		errString  string
		projectNil bool
	}{
		"malformed project": {
			response: `{
				"randomKey": "randomValue"
			}`,
			errString:  "Malformed",
			projectNil: true,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// Exact URL match
			httpmock.RegisterResponder(
				"GET",
				"http://localhost:8080/projects/1",
				httpmock.NewStringResponder(200, test.response),
			)

			client := NewKimaiClient("http://localhost:8080", "admin", "admin")

			project, err := client.GetProject(1)

			if test.errString != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.errString)
			}

			if test.projectNil {
				assert.Nil(t, project)
			} else {
				assert.NotNil(t, project)
			}

			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}

func TestCreateProject(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Exact URL match
	httpmock.RegisterResponder(
		"POST",
		"http://localhost:8080/projects",
		httpmock.NewStringResponder(200, `{
			"id": 1,
			"name": "Project 1",
			"customer": 1,
			"parentTitle": "",
			"comment": "",
			"visible": true,
			"metaFields": [],
			"start": "2006-01-02T15:04:05-0700",
			"end": "2006-01-02T15:04:05-0700"
		}`),
	)

	client := NewKimaiClient("http://localhost:8080", "admin", "admin")

	project := &Project{
		Name: "Project 1",
	}

	createdProject, err := client.CreateProject(project)

	assert.NoError(t, err)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())

	// Assert that a project is returned
	assert.Equal(t, "Project 1", createdProject.Name)
}
