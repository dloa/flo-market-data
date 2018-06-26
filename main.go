package main

import (
	"fmt"
	"github.com/oipwg/flo-market-data/fmd"
	"log"
	"net/http"
	"time"
)

func main() {
	err := fmd.InitDB("db/markets.db")
	if err != nil {
		log.Fatal(err)
	}

	apiKey := "apikey"
	secret := "secret"

	fmd.InitMRR(apiKey, secret)

	go fmd.WatchAndStoreForever(1 * time.Minute)

	fmdPrefix := ""
	http.Handle(fmdPrefix+"/", http.StripPrefix(fmdPrefix, fmd.ApiHandler))

	log.Println("Listening on port 41290")
	err = http.ListenAndServe(":41290", nil)
	if err != nil {
		log.Fatal("ListenAndServe failure: ", err)
		fmt.Printf("%v", err.Error())
	}
}
