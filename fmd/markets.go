package fmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	mrr "github.com/marccardinal/go-miningrigrentals-api"
)

var (
	mrrClient *mrr.Client
)

const cmcTickerUrl = "https://api.coinmarketcap.com/v1/ticker/%s/"

type cmcTicker []cmcTickerValue
type cmcTickerValue struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	VolUsd24h        string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	MaxSupply        string `json:"max_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

const bittrexGetMarketSummaryUrl = "https://bittrex.com/api/v1.1/public/getmarketsummary?market=btc-flo"

type bittrexGetMarketSummary struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Result  []bittrexGetMarketSummaryResult `json:"result"`
}
type bittrexGetMarketSummaryResult struct {
	MarketName     string  `json:"MarketName"`
	High           float64 `json:"High"`
	Low            float64 `json:"Low"`
	Volume         float64 `json:"Volume"`
	Last           float64 `json:"Last"`
	BaseVolume     float64 `json:"BaseVolume"`
	TimeStamp      string  `json:"TimeStamp"`
	Bid            float64 `json:"Bid"`
	Ask            float64 `json:"Ask"`
	OpenBuyOrders  int64   `json:"OpenBuyOrders"`
	OpenSellOrders int64   `json:"OpenSellOrders"`
	PrevDay        float64 `json:"PrevDay"`
	Created        string  `json:"Created"`
}

const niceHashLastUrl = "https://api.nicehash.com/api?method=stats.global.current"
const niceHash24hrUrl = "https://api.nicehash.com/api?method=stats.global.24h"
const niceHashScryptAlgo = 0

type niceHashApiSummary struct {
	Result niceHashApiResults `json:"result"`
	Method string             `json:"method"`
}

type niceHashApiResults struct {
	Stats []niceHashApiStats `json:"stats"`
}

type niceHashApiStats struct {
	Price string `json:"price"`
	Algo  int64  `json:"algo"`
	Speed string `json:"speed"`
}

type marketState struct {
	err               error
	time              time.Time
	cmcBtc            cmcTicker
	cmcFlo            cmcTicker
	cmcLtc            cmcTicker
	nhLast            niceHashApiSummary
	nh24Hr            niceHashApiSummary
	bittrexSummary    bittrexGetMarketSummary
	mrrRigsPriceInfo  mrr.RigListInfoPrice
	errCmcBtc         error
	errCmcFlo         error
	errCmcLtc         error
	errNhLast         error
	errNh24hr         error
	errBittrexSummary error
	errMrrRigs        error
}

func InitMRR(apiKey, secret string) {
	mrrClient = mrr.New(apiKey, secret)
}

// close(stop) to stop watching
func WatchMarkets(refreshRate time.Duration, stop <-chan struct{}, updates chan<- marketState) {
	defer close(updates)

	if fmdDB == nil {
		updates <- marketState{err: errors.New("must InitDB first")}
		return
	}

	run := true
	go func() {
		for run {
			updates <- refreshMarkets()
			time.Sleep(refreshRate)
		}
	}()

	<-stop
	run = false
	return
}

func refreshMarkets() (market marketState) {
	market.errCmcBtc = fetchJSON(fmt.Sprintf(cmcTickerUrl, "bitcoin"), &market.cmcBtc)
	market.errCmcFlo = fetchJSON(fmt.Sprintf(cmcTickerUrl, "florincoin"), &market.cmcFlo)
	market.errCmcFlo = fetchJSON(fmt.Sprintf(cmcTickerUrl, "litecoin"), &market.cmcLtc)
	market.errNhLast = fetchJSON(niceHashLastUrl, &market.nhLast)
	market.errNh24hr = fetchJSON(niceHash24hrUrl, &market.nh24Hr)
	market.errBittrexSummary = fetchJSON(bittrexGetMarketSummaryUrl, &market.bittrexSummary)

	_, rigListInfo, err := mrrClient.ListRigs("scrypt", 1)
	market.errMrrRigs = err
	if err == nil {
		market.mrrRigsPriceInfo = rigListInfo.Price
	}

	market.time = time.Now()
	return
}

func fetchJSON(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	resp.Body.Close()
	return err
}
