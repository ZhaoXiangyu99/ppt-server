package request

// Vue前端——>go server的请求
type ConversationParam struct {
	Data    string `json:"data" binding:"required"`
	Options Option `json:"options"`
}

type Option struct {
	ConversationId  string `json:"conversationId"`
	ParentMessageId string `json:"parentMessageId"`
}
