package logic

import (
	"context"
	"errors"
	"live/models/user_models"
	"live/utils/jwts"
	"live/utils/pwd"

	"live/auth_sever/auth_api/internal/svc"
	"live/auth_sever/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	var user user_models.UserModel
	err = l.svcCtx.DB.Take(&user, "id = ?", req.UserName).Error
	if err != nil {
		err = errors.New("用户名或密码错误")
		return nil, err
	}
	if !pwd.CheckPwd(user.Pwd, req.Password) {
		err = errors.New("用户名或密码错误")
		return nil, err
	}
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserID:   user.ID,
		Nickname: user.NickName,
		Role:     int(user.Role),
	}, l.svcCtx.Config.Auth.AccessSecret, int(l.svcCtx.Config.Auth.AccessExpire))
	if err != nil {
		logx.Error(err)
		err = errors.New("生成token失败")
		return nil, err
	}
	//context.WithValue(l.ctx, "userID", user.ID)
	//type Request struct {
	//	LogType int8   `json:"logType"`
	//	IP      string `json:"IP"`
	//	UserID  uint   `json:"userID"`
	//	Addr    string `json:"addr"`
	//	Level   string `json:"level"`
	//	Title   string `json:"title"`
	//	Content string `json:"content"`
	//	Service string `json:"service"`
	//}
	//requ := Request{
	//	LogType: 2,
	//	IP:      l.ctx.Value("clientIP").(string),
	//	UserID:  user.ID,
	//	Level:   "info",
	//	Title:   user.NickName + " 登录成功",
	//	Content: "登录成功",
	//	Service: l.svcCtx.Config.Name,
	//}
	//byteData, err := json.Marshal(requ)
	//if err != nil {
	//	logx.Errorf("json转化失败:%v", err)
	//}
	//
	//err = l.svcCtx.KqPusherClient.Push(string(byteData))
	//if err != nil {
	//	fmt.Printf("push error:%v", err)
	//}
	return &types.LoginResponse{Token: token}, nil
}
