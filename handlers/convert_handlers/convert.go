package convert_handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/attributeerror/currency-rates-service/database"
	"github.com/gin-gonic/gin"
)

var GetConvertCurrencyFromBase = func(r database.Database) func(c *gin.Context) {
	return func(c *gin.Context) {
		baseCode := c.Query("base")
		toCode := c.Query("to")
		amount := c.Query("amount")

		if baseCode == "" || toCode == "" || amount == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing one or more required parameters: base, to, amount"})
			return
		}

		baseAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error whilst parsing amount: %s", err.Error())})
			return
		}

		fmt.Printf("received request to convert %f %s to %s\n", baseAmount, baseCode, toCode)

		if baseCode != "EUR" {
			baseCodeToEuroRate, err := r.GetEuroRateForCurrency(baseCode) // 1 EUR = 1.09164 USD
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error whilst getting euro rate for currency %s: %s", baseCode, err.Error())})
				return
			}
			toCodeToEuroRate, err := r.GetEuroRateForCurrency(toCode) // 1 EUR = 0.87633 GBP
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error whilst getting euro rate for currency %s: %s", toCode, err.Error())})
				return
			}

			fmt.Printf("1 EUR = %f %s\n1 EUR = %f %s\n", *baseCodeToEuroRate, baseCode, *toCodeToEuroRate, toCode)

			baseToEurRate := 1 / *baseCodeToEuroRate // 1 USD = 0.91605 EUR

			baseAmountInEuros := baseAmount * baseToEurRate // 20 USD * 0.91605 EUR = 18.321 EUR

			convertedAmount := baseAmountInEuros * *toCodeToEuroRate // 18.321 EUR * 0.87633 GBP = 16.05524 GBP

			fmt.Printf("%f %s = %f %s\n", baseAmount, baseCode, convertedAmount, toCode)

			c.JSON(http.StatusOK, gin.H{"from": baseCode, "to": toCode, "baseAmount": baseAmount, "convertedAmount": convertedAmount})
			return
		}

		euroRate, err := r.GetEuroRateForCurrency(toCode) // 1 EUR = 0.87633 GBP
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error whilst getting euro rate for currency %s: %s", toCode, err.Error())})
			return
		}

		fmt.Printf("1 EUR = %f %s\n", *euroRate, toCode)

		convertedAmount := baseAmount * *euroRate // 20 EUR = 17.5266 GBP

		c.JSON(http.StatusOK, gin.H{"from": baseCode, "to": toCode, "baseAmount": baseAmount, "convertedAmount": convertedAmount})
	}
}
