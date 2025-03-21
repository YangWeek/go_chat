package logic

import (
	"context"
	"easy-chat/apps/user/models"
	"errors"
	"github.com/jinzhu/copier"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNotFound = errors.New("此用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	// todo: add your logic here and delete this line
	// 通过用户ID从数据库中查找用户信息
	userEntity, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		// 处理用户不存在的情况
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 将数据库中的用户实体复制到响应对象中
	var resp user.UserEntity
	if err := copier.Copy(&resp, userEntity); err != nil {
		return nil, err
	}

	// 返回用户信息
	return &user.GetUserInfoResp{
		User: &resp,
	}, nil
}
