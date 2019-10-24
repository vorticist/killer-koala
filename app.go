package main

import (
	"html/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vorticist/killer-koala/auth"
	"github.com/vorticist/killer-koala/routing"
	"github.com/vorticist/logger"
	mgo "gopkg.in/mgo.v2"
)

type App struct {
	Config           *AppConfig
	nonSecuredRoutes routing.Routes
	securedRoutes    routing.Routes
	appViews         []string
	session          *mgo.Session
	db               *mgo.Database
}

type AppConfig struct {
	MongoDBUrl  string
	MongoDBName string
	PublicKey   string
	PrivateKey  string
	Port        string
}

type Router interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

func NewAppWithConfig(config *AppConfig) App {
	var err error
	app := App{Config: config}
	app.session, err = mgo.Dial(app.Config.MongoDBUrl)
	if err != nil {
		panic(err)
	}
	app.session.SetMode(mgo.Monotonic, true)

	app.db = app.session.DB(app.Config.MongoDBName)
	return app
}

func (a *App) Serve() {
	defer a.session.Close()
	if a.Config == nil {
		logger.Error("not valid config values")
		return
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, ${error} \t| ${latency_human}\n", //"${time_rfc3339} ${id} ${short_file} ${line}",
	}))
	e.Use(middleware.Recover())

	mapRoutes(e, a.nonSecuredRoutes)

	auth.InitKeys(a.Config.PublicKey, a.Config.PrivateKey)

	r := e.Group("/api")

	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: auth.GetPrivateKey(),
	}))

	mapRoutes(r, a.securedRoutes)

	t := &routing.Template{
		Templates: template.Must(template.ParseFiles(a.appViews...)),
	}
	e.Renderer = t

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":" + a.Config.Port))
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
