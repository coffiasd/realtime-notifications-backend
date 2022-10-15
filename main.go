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

func Go(x func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Sprintf("panic %s\n", err))
			}
		}()
		x()
	}()
}

func main() {
	//listen blockchain data
	Go(func() {
		cron.Fetch()
	})

	httpServer := router.InitRouter()
	server := &http.Server{Addr: conf.App.HttpListen, Handler: httpServer}
	fmt.Println("server is running in", conf.App.HttpListen)
	server.ListenAndServe()
}
