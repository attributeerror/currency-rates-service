package convert_handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/attributeerror/currency-rates-service/database"
	"github.com/attributeerror/currency-rates-service/handlers/response_models"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	ErrMissingRequiredParameters = errors.New("missing one or more required parameters: base, to, amount")
)

var GetConvertCurrencyFromBase = func(r database.Database, sfGroup *singleflight.Group) func(c *gin.Context) {
	return func(c *gin.Context) {
		baseCode := c.DefaultQuery("base", "EUR")
		toCode := c.DefaultQuery("to", "")
		amount := c.DefaultQuery("amount", "")

		response, err, _ := sfGroup.Do(fmt.Sprintf("%s-%s-%s", baseCode, toCode, amount), func() (interface{}, error) {
			if baseCode == "" || toCode == "" || amount == "" {
				return nil, ErrMissingRequiredParameters
			}

			baseAmount, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				return nil, fmt.Errorf("error whilst parsing amount: %w", err)
			}

			fmt.Printf("received request to convert %f %s to %s\n", baseAmount, baseCode, toCode)

			if baseCode != "EUR" {
				baseCodeToEuroRate, err := r.GetEuroRateForCurrency(baseCode)
				if err != nil {
					return nil, fmt.Errorf("error whilst getting euro rate for currency %s: %w", baseCode, err)
				}
				toCodeToEuroRate, err := r.GetEuroRateForCurrency(toCode)
				if err != nil {
					return nil, fmt.Errorf("error whilst getting euro rate for currency %s: %w", toCode, err)
				}

				fmt.Printf("1 EUR = %f %s\n1 EUR = %f %s\n", *baseCodeToEuroRate, baseCode, *toCodeToEuroRate, toCode)

				baseToEurRate := 1 / *baseCodeToEuroRate
				baseAmountInEuros := baseAmount * baseToEurRate
				convertedAmount := baseAmountInEuros * *toCodeToEuroRate

				fmt.Printf("%f %s = %f %s\n", baseAmount, baseCode, convertedAmount, toCode)

				return &response_models.ConvertCurrencyResponse{
					BaseCode:        baseCode,
					ToCode:          toCode,
					BaseAmount:      baseAmount,
					ConvertedAmount: convertedAmount,
				}, nil
			}

			euroRate, err := r.GetEuroRateForCurrency(toCode)
			if err != nil {
				return nil, fmt.Errorf("error whilst getting euro rate for currency %s: %s", toCode, err)
			}

			fmt.Printf("1 EUR = %f %s\n", *euroRate, toCode)

			convertedAmount := baseAmount * *euroRate

			return &response_models.ConvertCurrencyResponse{
				BaseCode:        baseCode,
				ToCode:          toCode,
				BaseAmount:      baseAmount,
				ConvertedAmount: convertedAmount,
			}, nil
		})

		if err != nil {
			var statusCode int
			if errors.Is(err, ErrMissingRequiredParameters) {
				statusCode = http.StatusBadRequest
			} else {
				statusCode = http.StatusInternalServerError
			}

			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
			return
		}

		if responseModel, ok := response.(*response_models.ConvertCurrencyResponse); ok {
			c.JSON(http.StatusOK, responseModel)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unknown error occurred whilst parsing the response. Please try again later.",
		})
	}
}
