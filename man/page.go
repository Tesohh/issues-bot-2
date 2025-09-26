package man

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

type Page struct {
	ID    string
	Title string
	Func  PageMaker
}

type PageMaker func(*dg.Session, *dg.Interaction) ([]dg.MessageComponent, error)

var Pages = map[string]Page{
	"shorthand":                 Shorthand,
	"priorities-and-categories": PrioritiesAndCategories,
	"dependencies":              Dependencies,
	"tasks":                     Tasks,
}

func dePijpToBackticks(s string) string {
	return strings.ReplaceAll(s, "|", "`")
}

func text(s string, args ...any) dg.TextDisplay {
	return dg.TextDisplay{Content: dePijpToBackticks(fmt.Sprintf(s, args...))}
}

func media(urls ...string) dg.MediaGallery {
	items := []dg.MediaGalleryItem{}
	for _, url := range urls {
		items = append(items, dg.MediaGalleryItem{Media: dg.UnfurledMediaItem{URL: url}})
	}
	return dg.MediaGallery{
		Items: items,
	}
}
