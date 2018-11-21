package main

import (
	"flag"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"os/signal"
)

var (
	app       *echo.Echo
	apiGroup  *echo.Group
	db        *pg.DB
	jwtAuth   echo.MiddlewareFunc
	jwtSecret = "AqsA46R925YquUaLvu5mGJNj"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:8000", "add server")
	flag.Parse()

	db = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "Conghuy.315",
		Database: "webtracker",
	})

	migrate()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func() {
		<-interruptChan
		db.Close()
		os.Exit(0)
	}()

	app = echo.New()
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Recover())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 9}))
	jwtAuth = middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserClaims{},
		SigningKey: []byte(jwtSecret),
	})

	apiGroup = app.Group("/api", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().Header.Get("api-key") == "c3Lc8Jc3DjBpQQw4PUgBgzxb" {
				return next(ctx)
			}
			return echo.ErrUnauthorized
		}
	})

	installedRouters := []BaseRouter{
		&AuthRouter{},
	}
	for i := range installedRouters {
		installedRouters[i].Install()
	}


	app.Start(addr)

}
