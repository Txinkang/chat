package model

// Messages 消息模型
type Messages struct {
	ID        string      `json:"id"`
	RoomId    string      `json:"room_id"`
	SenderId  string      `json:"sender_id"`
	Type      string      `json:"type"`
	Content   interface{} `json:"content"`
	CreatedAt int64       `json:"created_at"`
}
