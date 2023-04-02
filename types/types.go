package types

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Session struct {
	Messages []Message `json:"messages"`
	ID       int    `json:"id"`
}