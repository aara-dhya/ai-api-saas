package ai

import "errors"

type OpenAIProvider struct {
	apiKey string
}

func NewOpenAIProvider(key string) *OpenAIProvider {
	return &OpenAIProvider{apiKey: key}
}

func (o *OpenAIProvider) Generate(req GenerateRequest) (*GenerateResponse, error) {
	return nil, errors.New("openai provider not implemented")
}
