package responses

type ErrResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Id        string `json:"id"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}

type MessagesResponse struct {
	Messages   []MessageResponse `json:"messages"`
	Pagination Pagination        `json:"pagination"`
}

type Pagination struct {
	CurrentPage   int `json:"currentPage"`
	PageSize      int `json:"pageSize"`
	TotalMessages int `json:"totalMessages"`
	TotalPages    int `json:"totalPages"`
}
