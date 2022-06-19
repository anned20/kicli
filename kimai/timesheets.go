package kimai

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (t *Timesheet) GetRawStart() string {
	return t.RawStart
}

func (t *Timesheet) GetRawEnd() string {
	return t.RawEnd
}

func (t *Timesheet) SetStart(start *time.Time) {
	t.Start = start
}

func (t *Timesheet) SetEnd(end *time.Time) {
	t.End = end
}

func (t *Timesheet) RealDuration() time.Duration {
	if t.End == nil {
		return time.Since(*t.Start).Truncate(time.Duration(time.Second))
	}

	return t.End.Sub(*t.Start).Truncate(time.Duration(time.Second))
}

func (t *Timesheet) BilledDuration() time.Duration {
	return time.Duration(t.Duration) * time.Second
}

func (t *ActiveTimesheet) GetRawStart() string {
	return t.RawStart
}

func (t *ActiveTimesheet) GetRawEnd() string {
	return t.RawEnd
}

func (t *ActiveTimesheet) SetStart(start *time.Time) {
	t.Start = start
}

func (t *ActiveTimesheet) SetEnd(end *time.Time) {
	t.End = end
}

func (t *ActiveTimesheet) RealDuration() time.Duration {
	if t.End == nil {
		return time.Since(*t.Start).Truncate(time.Duration(time.Second))
	}

	return t.End.Sub(*t.Start).Truncate(time.Duration(time.Second))
}

func (k *KimaiClient) CreateTimesheet(timesheet *Timesheet) (*Timesheet, error) {
	createTimesheet := &CreateTimesheet{
		ProjectID:  timesheet.ProjectID,
		ActivityID: timesheet.ActivityID,
	}

	now := time.Now()

	createTimesheet.Start = now.Format(dateFormat)

	response, err := k.api.post("/timesheets", createTimesheet)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var newTimesheet *Timesheet

	err = json.NewDecoder(response.Body).Decode(&newTimesheet)

	if err != nil {
		return nil, err
	}

	if newTimesheet.ID == 0 {
		return nil, errors.New("Malformed timesheet")
	}

	newTs, err := normalizeTimestamps(newTimesheet)

	if err != nil {
		return nil, err
	}

	newTimesheet, ok := newTs.(*Timesheet)

	if !ok {
		return nil, fmt.Errorf("Could not convert to Timesheet")
	}

	return newTimesheet, nil
}

func (k *KimaiClient) GetActiveTimesheet() (*ActiveTimesheet, error) {
	response, err := k.api.get("/timesheets/active")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var activeTimesheets []ActiveTimesheet

	err = json.NewDecoder(response.Body).Decode(&activeTimesheets)

	if err != nil {
		return nil, err
	}

	if len(activeTimesheets) == 0 {
		return nil, nil
	}

	activeTimesheet := &activeTimesheets[0]

	if activeTimesheet.ID == 0 {
		return nil, errors.New("Malformed timesheet")
	}

	activeTs, err := normalizeTimestamps(activeTimesheet)

	if err != nil {
		return nil, err
	}

	activeTimesheet, ok := activeTs.(*ActiveTimesheet)

	if !ok {
		return nil, fmt.Errorf("Could not convert to ActiveTimesheet")
	}

	return activeTimesheet, nil
}

func (k *KimaiClient) GetTimesheet(timesheetID int) (*Timesheet, error) {
	response, err := k.api.get(fmt.Sprintf("/timesheets/%d", timesheetID))

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var timesheet *Timesheet

	err = json.NewDecoder(response.Body).Decode(&timesheet)

	if err != nil {
		return nil, err
	}

	if timesheet.ID == 0 {
		return nil, errors.New("Malformed timesheet")
	}

	ts, err := normalizeTimestamps(timesheet)

	if err != nil {
		return nil, err
	}

	timesheet, ok := ts.(*Timesheet)

	if !ok {
		return nil, fmt.Errorf("Could not convert to Timesheet")
	}

	return timesheet, nil
}

func (k *KimaiClient) StopTimesheet(timesheetID int) (*Timesheet, error) {
	response, err := k.api.patch(fmt.Sprintf("/timesheets/%d/stop", timesheetID), nil)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var timesheet *Timesheet

	err = json.NewDecoder(response.Body).Decode(&timesheet)

	if err != nil {
		return nil, err
	}

	if timesheet.ID == 0 {
		return nil, errors.New("Malformed timesheet")
	}

	ts, err := normalizeTimestamps(timesheet)

	if err != nil {
		return nil, err
	}

	timesheet, ok := ts.(*Timesheet)

	if !ok {
		return nil, fmt.Errorf("Could not convert to Timesheet")
	}

	return timesheet, nil
}
