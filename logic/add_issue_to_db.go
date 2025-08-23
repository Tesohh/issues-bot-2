package logic

import (
	"issues/v2/db"
	"issues/v2/slash"
)

// what is required to be set on the issue by callers:
// Title, Tags, Kind, ProjectID, RecruiterUserID, AssigneeUsers, CategoryRoleID, PriorityRoleID
// btw for AssigneeUsers you just need to set an array of users...
//
// what needs to be done after calling this function
// - create discord thread, send message etc.
// - with a separate function, update the issue to add thread and message ID.
func AddIssueToDB(incompleteIssue *db.Issue) (*db.Issue, error) {
	issue := incompleteIssue // just an alias

	// set issue.Code to MAX(Code) + 1 else 0
	count, err := db.Issues.Where("project_id = ?", issue.ProjectID).Count(db.Ctx, "id")
	if err != nil {
		return nil, err
	}
	issue.Code = slash.Ptr(uint(count + 1))

	// add the issue to the database
	err = db.Issues.Create(db.Ctx, issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}
