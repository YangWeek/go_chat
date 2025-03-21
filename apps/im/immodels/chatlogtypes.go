package immodels

import (
	"easy-chat/pkg/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var DefaultChatLogLimit int64 = 100

// 聊天记录表
// MongoDB的聊天消息记录表‌，用于存储单条消息的完整元数据和内容
// 存储单条消息的完整内容与状态。
type ChatLog struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string             `bson:"conversationId"`
	SendId         string             `bson:"sendId"`
	RecvId         string             `bson:"recvId"`
	MsgFrom        int                `bson:"msgFrom"`
	ChatType       constants.ChatType `bson:"chatType"`
	MsgType        constants.MType    `bson:"msgType"`
	MsgContent     string             `bson:"msgContent"`
	SendTime       int64              `bson:"sendTime"`
	Status         int                `bson:"status"`
	ReadRecords    []byte             `bson:"readRecords"` // 记录该消息的已读信息

	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
