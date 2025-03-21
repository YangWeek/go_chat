package msgtransfer

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/bitmap"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MsgChatTransfer 处理聊天消息的转发。

type MsgChatTransfer struct {
	*baseMsgTransfer
}

// NewMsgChatTransfer 创建一个新的 MsgChatTransfer 实例。
//
// 该函数用于初始化并返回一个新的消息聊天转发器实例，它封装了基本消息转发的功能。
func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

// Consume 处理从消息队列中消费的聊天消息。
//
// 该方法从消息队列中获取的数据进行反序列化、记录日志，并将消息转发给目标用户。

func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key:", key, "value:", value)

	var (
		data  mq.MsgChatTransfer
		ctx   = context.Background()
		msgId = primitive.NewObjectID()
	)
	// 反序列化数据
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		MsgId:          data.MsgId,
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	// 得到默认已读的对象
	readRecord := bitmap.NewBitmap(0)
	// 发送消息者默认已读
	readRecord.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecord.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
