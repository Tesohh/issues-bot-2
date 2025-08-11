package db

type Guild struct {
	ID string `gorm:"primarykey"`

	DefaultPriorityRoleID string
	DefaultCategoryRoleID string
	NobodyRoleID          string
	DiscussionRoleID      string

	DefaultPriorityRole Role `gorm:"foreignKey:DefaultPriorityRoleID"`
	DefaultCategoryRole Role `gorm:"foreignKey:DefaultCategoryRoleID"`
	NobodyRole          Role `gorm:"foreignKey:NobodyRoleID"`
	DiscussionRole      Role `gorm:"foreignKey:DiscussionRoleID"`

	Roles    []Role
	Projects []Project
}
