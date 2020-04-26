package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"

	// "io/ioutil"
	"net/http"
	"time"
	dblayer "filestore-server/db"
	"filestore-server/util"
	cfg "filestore-server/config"
)



// PostSignupHandler 处理接口POST
func PostSignupHandler(c *gin.Context) {

username := c.Request.FormValue("username")
passwd := c.Request.FormValue("password")

if len(username) < 3 || len(passwd) < 5 {
c.JSON(http.StatusOK,gin.H{
	"msg":"Invalid parameter",
	"code": -1,
})
return
}
// 对密码进行加盐及取Sha1值加密
encPasswd := util.Sha1([]byte(passwd + cfg.PasswordSalt))
// 将用户信息注册到用户表中
suc := dblayer.UserSignUp(username, encPasswd)
if suc {
	c.JSON(http.StatusOK,gin.H{
		"msg":"Signup Success",
		"code": 0,
	})} else {
	c.JSON(http.StatusOK,gin.H{
		"msg":"Signup Failed",
		"code": -2,
	})}
}

// GetSignupHandler : GET获取注册
func GetSignupHandler(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// GetSignInHandler : 登录接口
func GetSignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// PostSignInHandler : 改写登陆接口POST method
func PostSignInHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	encPasswd := util.Sha1([]byte(password + cfg.PasswordSalt))

	// 1. 校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK,gin.H{
			"msg":"Login Failed",
			"code": -1,
		})
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		c.JSON(http.StatusOK,gin.H{
			"msg":"Login Failed",
			"code": -2,
		})
	}

	// 3. 登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Token    string
			Username string
			Location string
		}{
			Token: token,
			Username: username,
			Location: "/static/view/home.html",
		},
	}
	c.Data(http.StatusOK,"application/json",resp.JSONBytes())
}

// UserInfoHandler ： 查询用户信息
//noinspection GoUnresolvedReference
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	//	token := c.Request.FormValue("token")

	// // 2. 验证token是否有效
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusForbidden,
			gin.H{})
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}


// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {

	if len(token) != 40 {
		return false
	}

	// TODO: 判断token的时效性，是否过期
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[:8]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		return false
	}


	//user, err := dblayer.GetUserToken(token)
	//if err != nil {
	//	return false
	//}
	//if user.Token == token {
	//	return true
	//}
	//
	//

	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}