package logic

import "issues/v2/db"

func GetIssueRelationshipsOfKind(issue *db.Issue, kind db.RelationshipKind) (db.RelationshipsByDirection, error) {
	inbound, err := db.Relationships.
		Preload("FromIssue", nil).
		Preload("FromIssue.Tags", nil).
		Preload("FromIssue.PriorityRole", nil).
		Preload("FromIssue.CategoryRole", nil).
		Where("to_issue_id = ?", issue.ID).
		Where("kind = ?", kind).
		Find(db.Ctx)
	if err != nil {
		return db.RelationshipsByDirection{}, err
	}

	outbound, err := db.Relationships.
		Preload("ToIssue", nil).
		Preload("ToIssue.Tags", nil).
		Preload("ToIssue.PriorityRole", nil).
		Preload("ToIssue.CategoryRole", nil).
		Where("from_issue_id = ?", issue.ID).
		Where("kind = ?", kind).
		Find(db.Ctx)
	if err != nil {
		return db.RelationshipsByDirection{}, err
	}

	return db.RelationshipsByDirection{
		Inbound:  inbound,
		Outbound: outbound,
	}, nil
}
