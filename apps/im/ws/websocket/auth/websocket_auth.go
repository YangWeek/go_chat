package auth

import (
	"fmt"
	"net/http"
	"time"
)

// WebSocketAuth 是Authentication接口的一个实现，用于WebSocket的认证。
type WebSocketAuth struct{}

// WebSocketAuth 是Authentication接口的一个实现，用于WebSocket的认证。
// Authenticate 实现了Authentication接口的Authenticate方法。
// 对于WebSocket认证，假设总是成功的，因此返回true。
func (a *WebSocketAuth) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}

// 从url中
// ws://localhost:8080/ws?userId=123
// wss://api.example.com/chat?userId=alice
func (a *WebSocketAuth) UserId(r *http.Request) string {
	// 提取URL中的查询参数
	query := r.URL.Query()
	// 检查查询参数是否存在且包含"userId"
	if query != nil && query["userId"] != nil {
		// 如果存在"userId"参数，返回其值
		return fmt.Sprintf("%v", query["userId"])
	}

	// 如果不存在"userId"参数，生成并返回当前时间的Unix毫秒时间戳作为用户ID
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
