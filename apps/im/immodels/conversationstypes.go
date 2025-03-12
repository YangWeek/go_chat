package immodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 用户维度的会话列表信息
// 这个是读扩散
//
//	‌用户会话列表（Conversations表
type Conversations struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	UserId           string                   `bson:"userId"`           //用户唯一标识，用于快速检索用户所有会话
	ConversationList map[string]*Conversation `bson:"conversationList"` // 以map结构存储会话ID到会话元数据（Conversation）的映射，优化动态增删会话的性能

	// TODO: Fill your own fields  // 消息总数与序列号，用于未读数计算和版本控制
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
