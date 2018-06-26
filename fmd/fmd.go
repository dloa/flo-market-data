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
		var pVol, pBtcFlo, bVol, bBtcFlo, cBtcUsd, cFloUsd, cLtcUsd, vol, btc, usd, mrrLast10, mrrLast24hr float64

		if s.errPoloVol == nil {
			pVol, err = strconv.ParseFloat(s.poloVol.BtcFlo.Btc, 64)
			if err != nil {
				pVol = 0
			}
		} else {
			pVol = 0
		}

		if s.errPoloHistory == nil {
			pBtcFlo, err = strconv.ParseFloat(s.poloHistory[0].Rate, 64)
			if err != nil {
				pVol = 0
				pBtcFlo = 0
			}
		} else {
			pVol = 0
			pBtcFlo = 0
		}

		if s.errBittrexSummary == nil && s.bittrexSummary.Success == true {
			bVol = s.bittrexSummary.Result[0].BaseVolume
			bBtcFlo = s.bittrexSummary.Result[0].Last
		} else {
			bVol = 0
			bBtcFlo = 0
		}

		if s.errCmcBtc == nil {
			cBtcUsd, err = strconv.ParseFloat(s.cmcBtc[0].PriceUsd, 64)
			if err != nil {
				cBtcUsd = 0
			}
		} else {
			cBtcUsd = 0
		}

		if s.errCmcFlo == nil {
			cFloUsd, err = strconv.ParseFloat(s.cmcFlo[0].PriceUsd, 64)
			if err != nil {
				cFloUsd = 0
			}
		} else {
			cFloUsd = 0
		}

		if s.errCmcLtc == nil {
			cLtcUsd, err = strconv.ParseFloat(s.cmcLtc[0].PriceUsd, 64)
			if err != nil {
				cLtcUsd = 0
			}
		} else {
			cLtcUsd = 0
		}

		vol = pVol + bVol
		if vol == 0 {
			btc = 0
		} else {
			btc = truncate8((pBtcFlo*pVol + bBtcFlo*bVol) / vol)
		}
		usd = truncate8(cBtcUsd * btc)

		if s.errMrrRigs != nil {
			mrrLast10 = 0
		} else {
			mrrLast10 = s.mrrRigsPriceInfo.Last10
		}

		now := time.Now()
		dps, err := fetchDataPoint(now.Add(0 - time.Hour*24).Unix(), now.Unix(), 1440)
		if err != nil {
			mrrLast24hr = 0
		} else {
			mrrLast24hr = avgMrr(dps)
		}

		err = insertToDb(s.time.Unix(), pVol, pBtcFlo, bVol, bBtcFlo, cBtcUsd, cFloUsd, cLtcUsd, vol, btc, usd, mrrLast10, mrrLast24hr)
		if err != nil {
			log.Println("fmd: Dabatase insertion failed... ")
			log.Println(err)
		}
	}
}

func truncate8(f float64) float64 {
	return float64(int(f*1e8)) / 1e8
}

func avgMrr(dps []DataPoint) float64 {
	var avg float64 = 0
	var cnt int64 = 0
	for _, dp := range dps {
		if dp.MrrLast10 != 0 {
			cnt++
			avg += dp.MrrLast10
		}
	}
	if cnt > 100 {
		return avg / float64(cnt)
	}
	return 0
}
