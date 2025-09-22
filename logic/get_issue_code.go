package logic

import (
	"issues/v2/db"
)

func GetIssueCode(issue *db.Issue) (uint, error) {
	count, err := db.Issues.
		Where("project_id = ?", issue.ProjectID).
		Where("kind = ?", db.IssueKindNormal).
		Count(db.Ctx, "id")
	if err != nil {
		return 0, err
	}

	return uint(count + 1), nil
}
