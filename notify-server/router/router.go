package router

import (
	"net/http"
	"notify-server/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//cross options.
func Ops() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(cors.Default())
	router.Use(Ops())

	router.Use(gin.Recovery())

	sub := router.Group("/v1")
	{
		//return the list of subscription.
		sub.GET("/sub/sub-list", api.SubList)
		//user info.
		sub.GET("/sub/user-info", api.UserInfo)
		//subscribe events.
		sub.POST("/sub/save-sub", api.SaveSub)
		//save subscribed events.
		sub.POST("/sub/save-info", api.SaveInfo)
	}

	return router
}
