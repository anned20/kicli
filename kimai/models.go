package kimai

import "time"

// MetaField is the same for projects and tasks
type MetaField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Activity struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	ProjectID   int         `json:"project"`
	ParentTitle string      `json:"parentTitle"`
	Comment     string      `json:"comment"`
	Visible     bool        `json:"visible"`
	MetaFields  []MetaField `json:"metaFields"`
	Budget      int         `json:"budget"`
	TimeBudget  int         `json:"timeBudget"`
	BudgetType  string      `json:"budgetType"`
	Color       string      `json:"color"`
}

type Customer struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Number     string      `json:"number"`
	Comment    string      `json:"comment"`
	Visible    bool        `json:"visible"`
	Currency   string      `json:"currency"`
	MetaFields []MetaField `json:"metaFields"`
}

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Project struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	CustomerID  int         `json:"customer"`
	ParentTitle string      `json:"parentTitle"`
	Comment     string      `json:"comment"`
	Visible     bool        `json:"visible"`
	MetaFields  []MetaField `json:"metaFields"`

	RawStart string `json:"start"`
	RawEnd   string `json:"end"`
	Start    *time.Time
	End      *time.Time
}

type ActiveProject struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Customer    Customer    `json:"customer"`
	ParentTitle string      `json:"parentTitle"`
	Comment     string      `json:"comment"`
	Visible     bool        `json:"visible"`
	MetaFields  []MetaField `json:"metaFields"`

	RawStart string `json:"start"`
	RawEnd   string `json:"end"`
	Start    *time.Time
	End      *time.Time
}

type CreateProject struct {
	Name       string `json:"name"`
	CustomerID int    `json:"customer"`
	Visible    bool   `json:"visible"`
}

type Timesheet struct {
	ID           int      `json:"id"`
	ProjectID    int      `json:"project"`
	ActivityID   int      `json:"activity"`
	Description  string   `json:"description"`
	FixedRate    float64  `json:"fixedRate"`
	HourlyRate   float64  `json:"hourlyRate"`
	UserID       int      `json:"user"`
	Exported     bool     `json:"exported"`
	Billable     bool     `json:"billable"`
	Tags         []string `json:"tags"`
	Duration     int      `json:"duration"`
	Rate         float64  `json:"rate"`
	InternalRate float64  `json:"internalRate"`

	RawStart string `json:"begin"`
	RawEnd   string `json:"end"`
	Start    *time.Time
	End      *time.Time
}

type ActiveTimesheet struct {
	ID          int           `json:"id"`
	Project     ActiveProject `json:"project"`
	Activity    Activity      `json:"activity"`
	Description string        `json:"description"`
	FixedRate   float64       `json:"fixedRate"`
	HourlyRate  float64       `json:"hourlyRate"`
	Exported    bool          `json:"exported"`
	Billable    bool          `json:"billable"`
	Tags        []string      `json:"tags"`

	RawStart string `json:"begin"`
	RawEnd   string `json:"end"`
	Start    *time.Time
	End      *time.Time
}

type CreateTimesheet struct {
	ProjectID  int    `json:"project"`
	Start      string `json:"begin"`
	ActivityID int    `json:"activity"`
	Billable   bool   `json:"billable"`
}

type UserPreference struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type User struct {
	Language      string           `json:"language"`
	Timezone      string           `json:"timezone"`
	Preferences   []UserPreference `json:"preferences"`
	ID            int              `json:"id"`
	Username      string           `json:"username"`
	Enabled       bool             `json:"enabled"`
	Roles         []string         `json:"roles"`
	Alias         string           `json:"alias"`
	Title         string           `json:"title"`
	Avatar        string           `json:"avatar"`
	AccountNumber string           `json:"account_number"`
	Color         string           `json:"color"`
}

type Version struct {
	Version   string `json:"version"`
	VersionId int    `json:"versionId"`
	Candidate string `json:"candidate"`
	Semver    string `json:"semver"`
	Name      string `json:"name"`
	Copyright string `json:"copyright"`
}

// Because the API returns timestamps in the format 2006-01-02T15:04:05-0700,
// we need to convert them to time.Time objects.
// This interface declares the methods that are needed to convert the timestamps
// to time.Time objects.
type startAndEnd interface {
	GetRawStart() string
	GetRawEnd() string

	SetStart(start *time.Time)
	SetEnd(end *time.Time)
}
