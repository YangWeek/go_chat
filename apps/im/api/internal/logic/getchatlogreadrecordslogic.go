package logic

import (
	"context"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/user"
	"easy-chat/pkg/bitmap"
	"easy-chat/pkg/constants"

	"easy-chat/apps/im/api/internal/svc"
	"easy-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 根据消息id 连接这个grpc client 查询消息读取的情况
// 消息id 保存用户的客户端上
// 可以用于获取历史消息
func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	chatLogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		return nil, err
	}

	var (
		chatlog = chatLogs.List[0]
		reads   = []string{chatlog.SendId}
		unreads []string
		ids     []string
	)
	// 分别设置已读未读  和 消息类型
	switch constants.ChatType(chatlog.ChatType) {
	// 单聊的
	case constants.SingleChatType:
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId}
		} else {
			reads = append(reads, chatlog.RecvId)
		}
		ids = []string{chatlog.RecvId, chatlog.SendId}
	// 群聊的
	case constants.GroupChatType:
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatlog.RecvId,
		})
		if err != nil {
			return nil, err
		}

		// 使用位图技术高效存储群成员的已读状态
		aBitmap := bitmap.Load(chatlog.ReadRecords) // 位图
		for _, member := range groupUsers.List {
			ids = append(ids, member.UserId)
			if member.UserId == chatlog.SendId {
				continue
			}
			if aBitmap.IsSet(member.UserId) {
				reads = append(reads, member.UserId)
			} else {
				unreads = append(unreads, member.UserId)
			}
		}
	}

	userEntities, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	userEntitySet := make(map[string]*user.UserEntity, len(userEntities.User))
	for i, userEntity := range userEntities.User {
		userEntitySet[userEntity.Id] = userEntities.User[i]
	}

	// 设置手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}
	for i, unread := range unreads {
		if u := userEntitySet[unread]; u != nil {
			unreads[i] = u.Phone
		}
	}

	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
