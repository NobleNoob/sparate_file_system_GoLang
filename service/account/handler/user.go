package handler

import (
"context"
	"filestore-server/common"
	"filestore-server/service/account/proto"
	cfg "filestore-server/config"
	dblayer "filestore-server/db"
)


type User struct {}

func (u *User) Signin(ctx context.Context, req *proto.ReqSignin, res *proto.RespSignin) error {
	panic("implement me")
}

func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, res *proto.RespUserInfo) error {
	panic("implement me")
}

func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	panic("implement me")
}

func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	panic("implement me")
}

func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {

	username := req.Username
	password := req.Password
	encPasswd := util.Sha1([]byte(password + cfg.PasswordSalt))

	if len(username) < 3 || len(password) < 5 {
		res.Code = common.StatusParamInvalid
		res.Message = "Param Invalid"
		return nil
	}
	// 将用户信息注册到用户表中
	suc := dblayer.UserSignUp(username, encPasswd)
	if suc {
		res.Code = common.StatusOK
		res.Message = "Signup Successful"
	} else {
		res.Code = common.StatusRegisterFailed
		res.Message = "Register Failed"
	}
	return nil
}