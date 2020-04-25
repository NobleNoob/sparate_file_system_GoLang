package handler

import (
	"filestore-server/common"
	"filestore-server/util"
	"github.com/gin-gonic/gin"
	"net/http"
)
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")
		if len(username) < 3 || !IsTokenValid(token) {
			// Token 校验失败则提示
			c.Abort() //跳过gin handler 执行
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"Token Invalid",
				nil,
				)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}

}
