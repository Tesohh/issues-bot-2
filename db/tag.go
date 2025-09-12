package db

import (
	"fmt"
	"issues/v2/helper"
	"time"
)

type Tag struct {
	Name      string `gorm:"primaryKey"`
	ProjectID uint   `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Project Project
}

func (tag *Tag) Pretty(maxLen int) string {
	return fmt.Sprintf("`+%s`", helper.StrTrunc(tag.Name, maxLen))
}
