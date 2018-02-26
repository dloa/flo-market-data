package fmd

import (
	"errors"
	"log"
	"strconv"
	"time"
)

func WatchAndStoreForever(refreshRate time.Duration) error {
	if fmdDB == nil {
		return errors.New("must InitDB first")
	}

	stop := make(chan struct{})

	updates := make(chan marketState)
	go WatchMarkets(refreshRate, stop, updates)
	go storeMarketStateFromChannel(updates)

	<-stop
	return nil
}

func storeMarketStateFromChannel(states <-chan marketState) {
	for s := range states {
		var err error
		var pVol, pBtcFlo, bVol, bBtcFlo, cBtcUsd, cFloUsd, vol, btc, usd float64

		pVol, err = strconv.ParseFloat(s.poloVol.BtcFlo.Btc, 64)
		if err != nil {
			pVol = 0
		}
		pBtcFlo, err = strconv.ParseFloat(s.poloHistory[0].Rate, 64)
		if err != nil {
			pVol = 0
			pBtcFlo = 0
		}
		if s.bittrexSummary.Success == true {
			bVol = s.bittrexSummary.Result[0].BaseVolume
			bBtcFlo = s.bittrexSummary.Result[0].Last
		} else {
			bVol = 0
			bBtcFlo = 0
		}
		cBtcUsd, err = strconv.ParseFloat(s.cmcBtc[0].PriceUsd, 64)
		if err != nil {
			cBtcUsd = 0
		}
		cFloUsd, err = strconv.ParseFloat(s.cmcFlo[0].PriceUsd, 64)
		if err != nil {
			cFloUsd = 0
		}

		vol = pVol + bVol
		btc = truncate8((pBtcFlo*pVol + bBtcFlo*bVol) / vol)
		usd = truncate8(cBtcUsd * btc)

		err = insertToDb(s.time.Unix(), pVol, pBtcFlo, bVol, bBtcFlo, cBtcUsd, cFloUsd, vol, btc, usd)
		if err != nil {
			log.Println("fmd: Dabatase insertion failed... ")
			log.Println(err)
		}
	}
}

func truncate8(f float64) float64 {
	return float64(int(f*1e8)) / 1e8
}
