package immodels

import (
	"easy-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoDB的会话元数据表
// conversation记录会话列表【包含未读数，读取数据的节点】
type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	ConversationId string             `bson:"conversationId,omitempty"` // 会话唯一标识，规则化生成（如单聊user1_user2，群聊group_123）
	ChatType       constants.ChatType `bson:"chatType,omitempty"`
	//TargetId       string             `bson:"targetId,omitempty"`
	IsShow bool     `bson:"isShow,omitempty"`
	Total  int      `bson:"total,omitempty"`
	Seq    int64    `bson:"seq"`
	Msg    *ChatLog `bson:"msg,omitempty"` // 消息的单元

	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
