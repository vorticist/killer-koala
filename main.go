package main

import (
	"html/template"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vorticist/killer-koala/auth"
	"github.com/vorticist/killer-koala/routing"
	mgo "gopkg.in/mgo.v2"
)

var (
	nonSecuredRoutes = routing.Routes{}
	securedRoutes    = routing.Routes{}
	appRoutes        = routing.Routes{}
	appViews         = []string{}
)

type Router interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

func main() {
	dbHost := os.Getenv("MONGO_URL")
	session, err := mgo.Dial(dbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// db := session.DB(os.Getenv("MONGO_DB_NAME"))

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, ${error} \t| ${latency_human}\n", //"${time_rfc3339} ${id} ${short_file} ${line}",
	}))
	e.Use(middleware.Recover())

	// NOTE: Add all non secure routes to nonSecuredRoutes before this point
	mapRoutes(e, nonSecuredRoutes)

	auth.InitKeys(os.Getenv("PRIVATE_KEY"), os.Getenv("PUBLIC_KEY"))

	r := e.Group("/api")

	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: auth.GetPrivateKey(),
	}))

	// NOTE: Add all secure routes to securedRoutes before this point
	mapRoutes(r, securedRoutes)

	// NOTE: Add any and all view paths to appViews before this point
	t := &routing.Template{
		Templates: template.Must(template.ParseFiles(appViews...)),
	}
	e.Renderer = t
	mapRoutes(e, appRoutes)

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func mapRoutes(e Router, r routing.Routes) {
	for _, route := range r {
		switch route.Method {
		case routing.GET:
			e.GET(route.Pattern, route.HandleFunc)
			break
		case routing.POST:
			e.POST(route.Pattern, route.HandleFunc)
			break
		case routing.DELETE:
			e.DELETE(route.Pattern, route.HandleFunc)
			break
		case routing.PUT:
			e.PUT(route.Pattern, route.HandleFunc)
		}
	}
}
