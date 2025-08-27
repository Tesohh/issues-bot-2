package dataview

import "time"

// used for keeping information on the autolist display
type ProjectViewState struct {
	ProjectID uint

	CurrentPage int

	Filter IssueFilter
	Sorter IssueSorter

	LastMachineUpdate time.Time // last time we got updated from the backend
	LastUserUpdate    time.Time // last time user touched a knob
}

// maps MessageIDs to ProjectStates
var ProjectViewStates = map[string]ProjectViewState{}
