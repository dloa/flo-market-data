# flo-market-data

Calculates FLO market data from various markets, provides an API and pushes
updates to the block chain periodically.

## Install

flo-market-data is written in go.

It requires mattn/go-sqlite3 for database operations.

```
$ go get github.com/mattn/go-sqlite3
$ go get github.com/oipwg/flo-market-data
```

Optional: install sqlite3 locally!

*Ubuntu*: `sudo apt-get install sqlite3`

*OSX*: `brew install sqlite3`


## Running standalone

Navigate to the flo-market-data directory and run the program!

Remember to include all packages:

```
$ go run main.go
```


## API

Hit this URL with a `GET` request to see the recent market data:

```
http://127.0.0.1:41290/flo-market-data/v1/latest
```

You'll get a response like this:

```
  {
    "unixtime": 1518474167,
    "polo_vol": 2.47774214,
    "polo_btc_flo": 0.0000118,
    "bittrex_vol": 1.71902953,
    "bittrex_btc_flo": 0.0000118,
    "cmc_btc_usd": 8811.68,
    "cmc_flo_usd": 0.103969,
    "volume": 4.19677167,
    "weighted_btc": 0.00001179,
    "weighted_usd": 0.1038897
  }
```


# License

The MIT License

Copyright (c) 2013-2018 Flo Developers

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
