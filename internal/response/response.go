package response

type ConversationResponse struct {
	Token string `json:"token"`
}

// python server -> go server
type GPT struct {
	Type  string  `json:"type"`
	Data  GPTData `json:"data"`
	Error string  `json:"error"`
}

type GPTData struct {
	Id              string `json:"id"`
	Text            string `json:"text"`
	ConversationId  string `json:"conversationId"`
	ParentMessageId string `json:"parentMessageId"`
}
