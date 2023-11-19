package response_models

type (
	ConvertCurrencyResponse struct {
		ToCode          string
		BaseCode        string
		BaseAmount      float64
		ConvertedAmount float64
	}
)
