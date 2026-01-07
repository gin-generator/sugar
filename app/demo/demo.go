package main

import (
	"github.com/gin-generator/sugar/app/demo/route"
	"github.com/gin-generator/sugar/bootstrap"
	"github.com/gin-generator/sugar/middleware"
)

func main() {
	b := bootstrap.NewBootstrap(
		bootstrap.WithHttpMiddleware(
			middleware.Recovery(),
			middleware.Logger(),
			middleware.Cors(),
		), // add http global middleware
		bootstrap.WithHttpRouter(route.RegisterApi), // register http route handlers
	)
	b.Run()
}
