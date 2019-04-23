package bitrex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// marketSummaryResult holds the result of the response from the bitrex API, https://api.bittrex.com/api/v1.1/public/getmarketsummary for
// getting market summary
type marketSummaryResult struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  []marketSummary `json:"result"`
}

// marketSummary defines the rate for a single market
type marketSummary struct {
	MarketName string
	Rate       float64 `json:"Last"`
}

// DcrToUsd returns the value of the supplied amount in USD
func DcrToUsd(amountInDcr float64) (float64, error) {
	amountInBtc, err := dcrToBtc(amountInDcr)
	if err != nil {
		return 0, err
	}

	oneBtcInUsd, err := btcToUsd(1)
	if err != nil {
		return 0, err
	}

	return oneBtcInUsd * amountInBtc, nil
}

func dcrToBtc(amountInDcr float64) (float64, error) {
	summary, err := fetchMarketSummary("btc-dcr")
	if err != nil {
		return 0, err
	}
	return summary.Rate * amountInDcr, nil
}

func btcToUsd(btcAmount float64) (float64, error) {
	summary, err := fetchMarketSummary("usd-btc")
	if err != nil {
		return 0, err
	}
	return summary.Rate * btcAmount, nil
}

func fetchMarketSummary(marketName string) (*marketSummary, error) {
	url := fmt.Sprintf("https://api.bittrex.com/api/v1.1/public/getmarketsummary?market=%s", marketName)

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := marketSummaryResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return &result.Result[0], nil
}
