package fmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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

const poloVolFloUrl = "https://poloniex.com/public?command=return24hVolume"

type poloVolFlo struct {
	BtcFlo poloBtcFlo `json:"BTC_FLO"`
}
type poloBtcFlo struct {
	Btc string `json:"BTC"`
	Flo string `json:"FLO"`
}

const poloTradeHistoryUrl = "https://poloniex.com/public?command=returnTradeHistory&currencyPair=BTC_FLO"

type poloTradeHistory []poloTradeHistoryValue
type poloTradeHistoryValue struct {
	GlobalTradeId int64  `json:"globalTradeID"`
	TradeId       int64  `json:"tradeID"`
	Date          string `json:"date"`
	Type          string `json:"type"`
	Rate          string `json:"rate"`
	Amount        string `json:"amount"`
	Total         string `json:"total"`
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

type marketState struct {
	err               error
	time              time.Time
	cmcBtc            cmcTicker
	cmcFlo            cmcTicker
	poloVol           poloVolFlo
	poloHistory       poloTradeHistory
	bittrexSummary    bittrexGetMarketSummary
	errCmcBtc         error
	errCmcFlo         error
	errPoloVol        error
	errPoloHistory    error
	errBittrexSummary error
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
	market.errPoloVol = fetchJSON(poloVolFloUrl, &market.poloVol)
	market.errPoloHistory = fetchJSON(poloTradeHistoryUrl, &market.poloHistory)
	market.errBittrexSummary = fetchJSON(bittrexGetMarketSummaryUrl, &market.bittrexSummary)
	market.time = time.Now()
	return
}

func fetchJSON(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(result)
}
