package logic

import (
	"context"
	"easy-chat/apps/user/models"

	"easy-chat/apps/user/rpc/internal/svc"
	"easy-chat/apps/user/rpc/user"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line
	var (
		userEntities []*models.Users // 用于存储用户实体的切片
		err          error           // 存储可能出现的错误
	)

	// 根据不同的请求参数进行查询
	if in.Phone != "" {
		// 根据手机号查询用户
		userEntity, err := l.svcCtx.UsersModel.FindOneByPhoneNumber(l.ctx, in.Phone)
		if err == nil {
			// 如果查询成功，将用户实体添加到切片中
			userEntities = append(userEntities, userEntity)
		}
	} else if in.Name != "" {
		// 根据用户名查询用户列表
		userEntities, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		// 根据用户ID列表查询用户列表
		userEntities, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	// 处理查询过程中出现的错误
	if err != nil {
		return nil, err
	}

	// 将数据库中的用户实体复制到响应对象中
	var resp []*user.UserEntity
	if err := copier.Copy(&resp, userEntities); err != nil {
		return nil, err
	}

	// 返回查询结果
	return &user.FindUserResp{
		User: resp,
	}, nil
}
