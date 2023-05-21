package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	NBUUrl   = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?&json"
	USDIndex = 24
	EURIndex = 31
)

type Result struct {
	UAN         int
	Currency    string
	ResultValue float64
}

type NBUResponse []struct {
	R030         int     `json:"r030"`
	Txt          string  `json:"txt"`
	Rate         float64 `json:"rate"`
	Cc           string  `json:"cc"`
	ExchangeDate string  `json:"exchangedate"`
}

type ExchangeRatesKeeper struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
}

func NewExchangeRatesKeeper() *ExchangeRatesKeeper {
	var response NBUResponse

	resp, err := http.Get(NBUUrl)
	if err != nil {
		log.Fatal("an error occurred when get from nbu api ")
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("an error occurred when decode response from nbu api")
	}
	return &ExchangeRatesKeeper{
		USD: response[USDIndex].Rate,
		EUR: response[EURIndex].Rate,
	}
}

func (e *ExchangeRatesKeeper) ExchangeRatesGetter() {
	for {
		time.Sleep(5 * time.Minute)
		var response NBUResponse

		resp, err := http.Get(NBUUrl)
		if err != nil {
			log.Fatal("an error occurred when get from nbu api")
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal("an error occurred when decode response from nbu api")
		}
		e.USD = response[USDIndex].Rate
		e.EUR = response[EURIndex].Rate
	}
}

func (e *ExchangeRatesKeeper) currencyRate(c *gin.Context) {
	c.JSON(http.StatusOK, e)
}
