package slash

import (
	"fmt"
	"strings"
)

func MentionMany(ids []string, tag string, sep string) string {
	var mentions []string
	for _, id := range ids {
		mention := fmt.Sprintf("<%s%s>", tag, id)
		mentions = append(mentions, mention)
	}
	return strings.Join(mentions, sep)
}
