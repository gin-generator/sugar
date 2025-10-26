package main

import (
	"github.com/gin-generator/sugar/bootstrap"
	"github.com/gin-generator/sugar/middleware"
)

func main() {
	//v := bootstrap.NewConfig("env.yaml", "./etc")
	//fmt.Println(v.GetStringMap("database.mysql[]"))

	//mysql.NewMysql(v.Mysql)
	//mysql.DB.Default.Exec("show tables;")

	b := bootstrap.NewBootstrap(
		bootstrap.ServerHttp,
		bootstrap.WithHttpMiddleware(
			middleware.Recovery(),
			middleware.Logger(),
			middleware.Cors(),
		),
		bootstrap.WithHttpRouter(bootstrap.RegisterApiRoute),
	)
	b.Run()
}
