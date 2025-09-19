package logic

import "issues/v2/db"

func GetIssueRelationshipsOfKind(issue *db.Issue, kind db.RelationshipKind) ([]db.Relationship, error) {
	r, err := db.Relationships.
		Preload("ToIssue", nil).
		Preload("ToIssue.Tags", nil).
		Preload("ToIssue.PriorityRole", nil).
		Preload("ToIssue.CategoryRole", nil).
		Where("from_issue_id = ?", issue.ID).
		Where("kind = ?", kind).
		Find(db.Ctx)
	return r, err
}
