package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("home/index.html"))
}
