package models

import (
	"fmt"
	"time"

	"github.com/go-mods/initials"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/lib/pq"
)

type Post struct {
	ID          string `gorm:"primaryKey;default:cuid()"`
	CompanyName string
	Location    string
	Tags        pq.StringArray `gorm:"type:text[]"`
	Thumbnail   string
	Title       string
	Description string
	PublishedAt time.Time
	CreatedAt   time.Time
}

func (p Post) GetInitials() string {
	return initials.GetInitials(p.CompanyName)
}

func (p Post) TimeSinceCreated() string {
	timespan := time.Since(p.PublishedAt).Hours()
	switch {
	case timespan < 24:
		return "Today"
	case timespan < 48:
		return "Yesterday"
	case timespan < (24 * 30):
		return fmt.Sprintf("%d %s", int(timespan/24), "days ago")
	default:
		return fmt.Sprintf("%d %s", int(timespan/(24*30)), "months ago")
	}
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func (p Post) GetDescription() string {
	return string(mdToHTML([]byte(p.Description)))
}
