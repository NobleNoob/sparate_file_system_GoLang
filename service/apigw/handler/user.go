package handler

import (
	"context"
	"filestore-server/common"
	"filestore-server/util"
	"github.com/gin-gonic/gin"
	micro "github.com/micro/go-micro"
	"log"
	"net/http"

	cmn "filestore-server/common"
	userProto "filestore-server/service/account/proto"
	//_ "filestore-server/service/download/proto"
	//_ "filestore-server/service/upload/proto"
)

var (
	userCli userProto.UserService
	//upCli   upProto.UploadService
	//dlCli   dlProto.DownloadService
)

func init() {
	service := micro.NewService()

	// 初始化， 解析命令行参数等
	service.Init()

	// 初始化一个account服务的客户端
	userCli = userProto.NewUserService("go.micro.service.user", service.Client())
	//// 初始化一个upload服务的客户端
	//upCli = upProto.NewUploadService("go.micro.service.upload", service.Client())
	//// 初始化一个download服务的客户端
	//dlCli = dlProto.NewDownloadService("go.micro.service.download", service.Client())
}

// PostSignupHandler 处理接口POST
func PostSignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	ret,err := userCli.Signup(context.TODO(),&userProto.ReqSignup{
		Username: username,
		Password: passwd,
	})

	if err !=nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code": ret.Code,
		"message": ret.Message,
	})

}

// GetSignupHandler : GET获取注册
func GetSignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// SigninHandler : 响应登录页面
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// DoSigninHandler : 处理登录post请求
func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	rpcResp, err := userCli.Signin(context.TODO(), &userProto.ReqSignin{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if rpcResp.Code != cmn.StatusOK {
		c.JSON(200, gin.H{
			"msg":  "登录失败",
			"code": rpcResp.Code,
		})
		return
	}

	// 动态获取上传入口地址
	//upEntryResp, err := upCli.UploadEntry(context.TODO(), &upProto.ReqEntry{})
	//if err != nil {
	//	log.Println(err.Error())
	//} else if upEntryResp.Code != cmn.StatusOK {
	//	log.Println(upEntryResp.Message)
	//}
	//
	//// 动态获取下载入口地址
	//dlEntryResp, err := dlCli.DownloadEntry(context.TODO(), &dlProto.ReqEntry{})
	//if err != nil {
	//	log.Println(err.Error())
	//} else if dlEntryResp.Code != cmn.StatusOK {
	//	log.Println(dlEntryResp.Message)
	//}

	// 登录成功，返回用户信息
	cliResp := util.RespMsg{
		Code: int(common.StatusOK),
		Msg:  "登录成功",
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location:      "/static/view/home.html",
			Username:      username,
			Token:         rpcResp.Token,
			//UploadEntry:   upEntryResp.Entry,
			//DownloadEntry: dlEntryResp.Entry,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")

	resp, err := userCli.UserInfo(context.TODO(), &userProto.ReqUserInfo{
		Username: username,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 3. 组装并且响应用户数据
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username": username,
			"SignupAt": resp.SignupAt,
			// TODO: 完善其他字段信息
			"LastActive": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}