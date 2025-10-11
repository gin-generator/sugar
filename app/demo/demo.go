package main

import (
	"github.com/gin-generator/sugar/bootstrap"
)

func main() {
	//v := bootstrap.NewConfig("env.yaml", "./etc")
	//fmt.Println(v.GetStringMap("database.mysql[]"))

	//mysql.NewMysql(v.Mysql)
	//mysql.DB.Default.Exec("show tables;")

	b := bootstrap.NewBootstrap(bootstrap.ServerHttp)
	b.RunHttp()
}
