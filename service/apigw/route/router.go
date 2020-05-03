package route

import (
	"filestore-server/service/apigw/handler"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Static("/static/","./static")
	router.GET("/user/signup",handler.GetSignupHandler)
	router.POST("/user/signup",handler.PostSignupHandler)
	return router
}