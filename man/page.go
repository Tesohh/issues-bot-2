package man

import (
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
	"shorthand": Shorthand,
}

func dePijpToBackticks(s string) string {
	return strings.ReplaceAll(s, "|", "`")
}

func text(s string) dg.TextDisplay {
	return dg.TextDisplay{Content: dePijpToBackticks(s)}
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
