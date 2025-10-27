package main

import (
	_middleware "github.com/gin-generator/sugar/app/demo/middleware"
	"github.com/gin-generator/sugar/app/demo/route"
	"github.com/gin-generator/sugar/bootstrap"
	"github.com/gin-generator/sugar/middleware"
)

func main() {
	//v := bootstrap.NewConfig("env.yaml", "./etc")
	//fmt.Println(v.GetStringMap("database.mysql[]"))

	//mysql.NewMysql(v.Mysql)
	//mysql.DB.Default.Exec("show tables;")

	b := bootstrap.NewBootstrap(
		bootstrap.ServerHttp, // create a new http bootstrap instance
		bootstrap.WithHttpMiddleware(
			middleware.Recovery(),
			middleware.Logger(),
			middleware.Cors(),
			_middleware.Auth(), // add http server middleware
		), // add http global middleware
		bootstrap.WithHttpRouter(route.RegisterApi), // register http route handlers
	)
	b.Run()
}
