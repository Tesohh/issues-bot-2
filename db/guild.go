package db

type Guild struct {
	ID string `gorm:"primarykey"`

	GenericCategoryRoleID string
	FeatCategoryRoleID    string
	FixCategoryRoleID     string
	ChoreCategoryRoleID   string

	LowPriorityRoleID       string
	NormalPriorityRoleID    string
	ImportantPriorityRoleID string
	CriticalPriorityRoleID  string

	NobodyRoleID     string
	DiscussionRoleID string

	GenericCategoryRole Role `gorm:"foreignKey:GenericCategoryRoleID"`
	FeatCategoryRole    Role `gorm:"foreignKey:FeatCategoryRoleID"`
	FixCategoryRole     Role `gorm:"foreignKey:FixCategoryRoleID"`
	ChoreCategoryRole   Role `gorm:"foreignKey:ChoreCategoryRoleID"`

	LowPriorityRole       Role `gorm:"foreignKey:LowPriorityRoleID"`
	NormalPriorityRole    Role `gorm:"foreignKey:NormalPriorityRoleID"`
	ImportantPriorityRole Role `gorm:"foreignKey:ImportantPriorityRoleID"`
	CriticalPriorityRole  Role `gorm:"foreignKey:CriticalPriorityRoleID"`

	NobodyRole     Role `gorm:"foreignKey:NobodyRoleID"`
	DiscussionRole Role `gorm:"foreignKey:DiscussionRoleID"`

	Roles    []Role
	Projects []Project
}
