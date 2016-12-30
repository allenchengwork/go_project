package main

import (
	"gopkg.in/iris-contrib/middleware.v5/logger"
	"gopkg.in/iris-contrib/middleware.v5/recovery"
	"gopkg.in/kataras/go-template.v0/html"
	"gopkg.in/kataras/iris.v5"
	"maplebox/conf"
	"maplebox/controller"
	"maplebox/logs"
	"net/http"
)

var (
	router *iris.Framework
	server *http.Server
	log    *logs.AppLogger = logs.AppLog
)

func main() {
	defer recoverApp()

	initConfig()

	initLog()

	initRouter()

	log.Info("Init App End")
	startServer()
}

func initConfig() {
	conf.InitConfig()
}

func initLog() {
	logs.InitLog()
}

func initRouter() {
	var isDev bool
	switch conf.AppConfig.AppMode {
	case "dev":
		isDev = true
	case "test":
		isDev = false
	case "prod":
		isDev = false
	default:
		isDev = true
	}

	router = iris.New(iris.Configuration{
		ReadTimeout:       conf.GetReadTimeout(),
		WriteTimeout:      conf.GetWriteTimeout(),
		DisablePathEscape: true,
		IsDevelopment:     isDev,
		DisableBanner:     true,
		LoggerOut:         logs.RouteLog.Logger.Out,
	})

	router.Use(logger.New())
	router.Use(recovery.Handler)

	router.Static("/static", conf.GetStaticPath(), 1)
	router.UseTemplate(html.New(html.Config{
		Layout: "layout.html",
	})).Directory(conf.GetViewsPath(), ".html")

	router.Get("/", func(ctx *iris.Context) {
		if err := ctx.Render("index.htm", iris.Map{
			"title": "Posts",
		}); err != nil {
			ctx.EmitError(iris.StatusServiceUnavailable)
		}
	})

	router.Get("/admin", func(ctx *iris.Context) {
		if err := ctx.Render("admin/index.html", iris.Map{
			"title": "Gin",
		}); err != nil {
			ctx.EmitError(iris.StatusServiceUnavailable)
		}
	})

	router.OnError(iris.StatusServiceUnavailable, func(ctx *iris.Context) {
		ctx.WriteString("ServiceUnavailable")
	})

	router.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.WriteString("InternalServerError")
	})

	authorized := router.Party("/api")
	authorized.UseFunc(authRequired())
	{
		var employeeController controller.EmployeeController
		authorized.Get("/employee", employeeController.EmployeeList)
	}
}

func recoverApp() {

}

func startServer() {
	router.Listen(conf.GetAddr())
}

func authRequired() iris.HandlerFunc {
	return func(c *iris.Context) {
		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
	}
}
