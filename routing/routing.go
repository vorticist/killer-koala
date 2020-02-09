package routing

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Route struct {
	Name       string
	HTTPVerb   string
	Pattern    string
	HandleFunc func(echo.Context) error
}

type Routes []Route

type Handler interface {
	Routes() Routes
	SecuredRoutes() Routes
}

type ViewHandler interface {
	Handler
	Views() []string
}

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
