package ai

import "errors"

type Router struct {
	providers map[string]Provider
}

func NewRouter() *Router {
	return &Router{
		providers: make(map[string]Provider),
	}
}

func (r *Router) Register(model string, provider Provider) {
	r.providers[model] = provider
}

func (r *Router) Generate(req GenerateRequest) (*GenerateResponse, error) {

	provider, ok := r.providers[req.Model]

	if !ok {
		return nil, errors.New("model not supported")
	}

	return provider.Generate(req)
}
