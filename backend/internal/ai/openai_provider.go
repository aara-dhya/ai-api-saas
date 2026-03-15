package ai

type OpenAIProvider struct {
	apiKey string
}

func NewOpenAIProvider(key string) *OpenAIProvider {
	return &OpenAIProvider{apiKey: key}
}

func (o *OpenAIProvider) Generate(req GenerateRequest) (*GenerateResponse, error) {
	// call OpenAI chat completions
}
