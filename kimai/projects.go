package kimai

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (p *Project) GetRawStart() string {
	return p.RawStart
}

func (p *Project) GetRawEnd() string {
	return p.RawEnd
}

func (p *Project) SetStart(start *time.Time) {
	p.Start = start
}

func (p *Project) SetEnd(end *time.Time) {
	p.End = end
}

func (p *ActiveProject) GetRawStart() string {
	return p.RawStart
}

func (p *ActiveProject) GetRawEnd() string {
	return p.RawEnd
}

func (p *ActiveProject) SetStart(start *time.Time) {
	p.Start = start
}

func (p *ActiveProject) SetEnd(end *time.Time) {
	p.End = end
}

func (k *KimaiClient) GetProjects() ([]Project, error) {
	response, err := k.api.get("/projects")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var projects []Project

	err = json.NewDecoder(response.Body).Decode(&projects)

	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return []Project{}, nil
	}

	if projects[0].ID == 0 {
		return nil, errors.New("Malformed projects")
	}

	// Convert the timestamps to time.Time
	for _, project := range projects {
		_p, err := normalizeTimestamps(&project)

		if err != nil {
			return nil, err
		}

		p, ok := _p.(*Project)

		if !ok {
			return nil, fmt.Errorf("Could not convert to Project")
		}

		project = *p
	}

	return projects, nil
}

func (k *KimaiClient) GetProject(id int) (*Project, error) {
	response, err := k.api.get("/projects/" + fmt.Sprintf("%d", id))

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var project *Project

	err = json.NewDecoder(response.Body).Decode(&project)

	if err != nil {
		return nil, err
	}

	if project.ID == 0 {
		return nil, errors.New("Malformed projects")
	}

	p, err := normalizeTimestamps(project)

	if err != nil {
		return nil, err
	}

	project, ok := p.(*Project)

	if !ok {
		return nil, fmt.Errorf("Could not convert to Project")
	}

	return project, nil
}

func (k *KimaiClient) CreateProject(project *Project) (*Project, error) {
	createProject := &CreateProject{
		Name:             project.Name,
		CustomerID:       project.CustomerID,
		Visible:          true,
		Billable:         true,
		GlobalActivities: true,
	}

	response, err := k.api.post("/projects", createProject)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var newProject *Project

	err = json.NewDecoder(response.Body).Decode(&newProject)

	if err != nil {
		return nil, err
	}

	newP, err := normalizeTimestamps(newProject)

	if err != nil {
		return nil, err
	}

	newProject, ok := newP.(*Project)

	if !ok {
		return nil, fmt.Errorf("Could not convert to Project")
	}

	return newProject, nil
}
