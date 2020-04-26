package handler

import (
"context"
	"filestore-server/common"
	"filestore-server/service/account/proto"
)
type User struct {

}

func Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {

	username := req.Username
	password := req.Password

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	if len(username) < 3 || len(password) <5 {
		res.Code = common.StatusParamInvalid
		res.Message = "Param Invalid"
		return nil
	}

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
