package billing

var pricePerToken = map[string]float64{
	"llama-3.1-8b-instant": 0.0000002,
}

func CalculateCost(model string, tokens int) float64 {

	price, ok := pricePerToken[model]
	if !ok {
		return 0
	}

	return float64(tokens) * price
}
