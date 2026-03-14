package ai

type Provider interface {
	Generate(req GenerateRequest) (*GenerateResponse, error)
}
