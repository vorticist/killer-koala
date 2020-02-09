package koala

import (
	"context"
	"html/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vorticist/killer-koala/auth"
	"github.com/vorticist/killer-koala/routing"
	"github.com/vorticist/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Config           *AppConfig
	nonSecuredRoutes routing.Routes
	securedRoutes    routing.Routes
	appViews         []string
	client           *mongo.Client
	db               *mongo.Database
	middlewares      []echo.MiddlewareFunc
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

	if len(app.Config.MongoDBUrl) > 0 && len(config.MongoDBName) > 0 {
		clientOptions := options.Client().ApplyURI(app.Config.MongoDBUrl)
		app.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			panic(err)
		}
	}
	return app
}

func (a *App) AddHandler(handler routing.Handler) {
	a.securedRoutes = append(a.securedRoutes, handler.SecuredRoutes()...)
	a.nonSecuredRoutes = append(a.nonSecuredRoutes, handler.Routes()...)
}

func (a *App) AddViewHandler(viewHandler routing.ViewHandler) {
	a.AddHandler(viewHandler)
	a.appViews = append(a.appViews, viewHandler.Views()...)
}

func (a *App) Database() *mongo.Database {
	if a.db == nil {
		a.db = a.client.Database(a.Config.MongoDBName)
	}
	return a.db
}

func (a *App) AddMiddleware(mf echo.MiddlewareFunc) {
	a.middlewares = append(a.middlewares, mf)
}

func (a *App) Serve() {
	if a.Config == nil {
		logger.Error("not valid config values")
		return
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, ${error} \t| ${latency_human}\n", //"${time_rfc3339} ${id} ${short_file} ${line}",
	}))
	e.Use(a.middlewares...)
	e.Use(middleware.Recover())

	mapRoutes(e, a.nonSecuredRoutes)

	auth.InitKeys(a.Config.PublicKey, a.Config.PrivateKey)

	r := e.Group("/api")

	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: auth.GetPrivateKey(),
	}))

	mapRoutes(r, a.securedRoutes)

	if len(a.appViews) > 0 {
		t := &routing.Template{
			Templates: template.Must(template.ParseFiles(a.appViews...)),
		}
		e.Renderer = t
	}

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":" + a.Config.Port))
}

func mapRoutes(e Router, r routing.Routes) {
	for _, route := range r {
		switch route.HTTPVerb {
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
