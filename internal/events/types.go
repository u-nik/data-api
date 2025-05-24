package events

type BaseEvent struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}
