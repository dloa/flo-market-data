package fmd

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var (
	ApiHandler http.Handler
)

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/last", apiGetLast)
	mux.HandleFunc("/range", apiGetRange)

	ApiHandler = cors(mux)
}

type DataPoint struct {
	Unixtime      float64 `json:"unixtime"`
	PoloVol       float64 `json:"polo_vol"`
	PoloBtcFlo    float64 `json:"polo_btc_flo"`
	BittrexVol    float64 `json:"bittrex_vol"`
	BittrexBtcFlo float64 `json:"bittrex_btc_flo"`
	CmcBtcUsd     float64 `json:"cmc_btc_usd"`
	CmcLtcUsd     float64 `json:"cmc_ltc_usd"`
	CmcFloUsd     float64 `json:"cmc_flo_usd"`
	Volume        float64 `json:"volume"`
	WeightedBtc   float64 `json:"weighted_btc"`
	WeightedUsd   float64 `json:"weighted_usd"`
	MrrLast10     float64 `json:"mrr_last_10"`
	MrrLast24hr   float64 `json:"mrr_last_24hr"`
	NhLast        float64 `json:"nh_last"`
	Nh24hr        float64 `json:"nh_24hr"`
}

func apiGetLast(w http.ResponseWriter, _ *http.Request) {
	res, err := fetchDataPoint(0, time.Now().Unix(), 1)
	if err != nil {
		http.Error(w, "Unable to obtain results", http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(res[0])
	if err != nil {
		http.Error(w, "Unable to encode results", http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func apiGetRange(w http.ResponseWriter, r *http.Request) {
	var from, to, limit int64
	var err error

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse parameters", http.StatusBadRequest)
	}
	if from, err = strconv.ParseInt(r.Form.Get("from"), 10, 64); err != nil {
		from = 0
	}
	if to, err = strconv.ParseInt(r.Form.Get("to"), 10, 64); err != nil {
		to = time.Now().Unix()
	}
	if limit, err = strconv.ParseInt(r.Form.Get("limit"), 10, 64); err != nil {
		limit = 10
	}

	res, err := fetchDataPoint(from, to, limit)
	if err != nil {
		http.Error(w, "Unable to obtain results", http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Unable to encode results", http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET")
		h.ServeHTTP(w, r)
	})
}
