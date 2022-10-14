package main

import (
	"fmt"
	"net/http"
	"notify-server/config"
	"notify-server/cron"
	"notify-server/middleware"
	"notify-server/router"
)

var (
	conf config.Conf
)

// init log , config,  redis server.
func init() {
	//load configs
	config.Parse("./conf/config.toml")
	conf = config.ParseConfig

	//init log
	middleware.InitLog()

	// init redis
	middleware.InitRedis()
}

func main() {
	//listen blockchain data
	go cron.Fetch()

	httpServer := router.InitRouter()
	server := &http.Server{Addr: conf.App.HttpListen, Handler: httpServer}
	fmt.Println("server is running in", conf.App.HttpListen)
	server.ListenAndServe()
}
