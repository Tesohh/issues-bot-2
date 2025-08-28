package db

import "time"

type RoleKind string

const (
	RoleKindPriority   RoleKind = "priority"
	RoleKindCategory   RoleKind = "category"
	RoleKindDiscussion RoleKind = "discussion"
	RoleKindNobody     RoleKind = "nobody"
)

type Role struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time

	Kind  RoleKind `gorm:"check:kind in ('priority', 'category', 'discussion', 'nobody')"`
	Key   string
	Emoji string

	GuildID string
}
