package components

import (
	"fmt"
	"vauntly/models"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func PostSearchResults(posts []models.Post, selectedTags string, nextPage int) g.Node {
	nextPgUrl := fmt.Sprint("/partials/posts/search/", nextPage)
	if selectedTags != "" {
		nextPgUrl = fmt.Sprint("/partials/posts/search/", nextPage, "?tags=", selectedTags)
	}

	return h.Div(
		g.Group(g.Map(posts, func(post models.Post) g.Node {
			return ResultsRow(post)
		})),
		h.Div(
			g.Attr("hx-post", nextPgUrl),
			g.Attr("hx-swap", "outerHTML"),
			g.Attr("hx-trigger", "revealed"),
			h.ID("nextPageLoaderId_"+fmt.Sprint(nextPage)),
			h.Class("htmx-indicator flex flex-col gap-2 items-center justify-center"),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
			h.Div(
				h.Class("bg-orange-200 dark:bg-slate-800 rounded-md h-36 w-full"),
			),
		),
	)
}
