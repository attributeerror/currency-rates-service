package response_models

type (
	ConvertCurrencyResponse struct {
		ToCode          string  `json:"to,omitempty"`
		BaseCode        string  `json:"from,omitempty"`
		BaseAmount      float64 `json:"baseAmount,omitempty"`
		ConvertedAmount float64 `json:"convertedAmount,omitempty"`
	}
)
