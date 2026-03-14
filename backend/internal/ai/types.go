package ai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GenerateRequest struct {
	Prompt string
	Model  string
}

type GenerateResponse struct {
	Model  string
	Tokens int
	Raw    []byte
}
