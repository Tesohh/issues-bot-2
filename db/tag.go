package db

import (
	"fmt"
	"issues/v2/helper"
	"strings"
	"time"
)

func ParseTags(raw string) []string {
	tags := []string{}
	for rawTag := range strings.SplitSeq(raw, ",") {
		trim := strings.Trim(rawTag, " +")
		if len(trim) > 0 {
			tags = append(tags, trim)
		}
	}
	return tags
}

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
