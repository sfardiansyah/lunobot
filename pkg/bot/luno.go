package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// PriceCandle ...
type PriceCandle struct {
	Candles []Candle `json:"candles"`
}

// Candle ...
type Candle struct {
	High float64 `json:"high,string"`
	Low  float64 `json:"low,string"`
}

// Ticker ...
type Ticker struct {
	Volume float64 `json:"rolling_24_hour_volume,string"`
}

// PairResponse ...
type PairResponse struct {
	Pairs []Pair `json:"availablePairs"`
}

// Pair ...
type Pair struct {
	BaseCode    string  `json:"baseCode"`
	CounterCode string  `json:"counterCode"`
	Price       float64 `json:"price,string"`
}

const (
	lunoAPIURL         = "https://api.mybitx.com/api/1/ticker"
	lunoPriceChartURL  = "https://www.luno.com/ajax/1/price_chart"
	lunoChartCandleURL = "https://www.luno.com/ajax/1/charts_candles"
)

func getInfo() string {
	idTime, _ := time.LoadLocation("Asia/Jakarta")
	date := time.Now().In(idTime).Format(time.RFC1123)

	btc := getPrice("XBTIDR")
	eth := getPrice("ETHIDR")
	btcH, btcL := getHiLo("XBTIDR")
	ethH, ethL := getHiLo("ETHIDR")

	return fmt.Sprintf(fileReader("assets/info.txt"), date, btc, btcH, btcL, eth, ethH, ethL)
}

func getHiLo(pair string) (string, string) {
	p := new(PriceCandle)
	date := time.Now().Add(-24 * time.Hour).Unix()

	r, err := http.Get(lunoChartCandleURL + "?pair=" + pair + "&since=" + strconv.FormatInt(date, 10))
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewDecoder(r.Body).Decode(p); err != nil {
		log.Fatal(err)
	}

	var hi, lo float64 = 0, p.Candles[0].Low
	for _, candle := range p.Candles {
		if candle.High > hi {
			hi = candle.High
		}
		if candle.Low < lo {
			lo = candle.Low
		}
	}

	return humanize.Commaf(hi), humanize.Commaf(lo)
}

func getPrice(pair string) string {
	p := new(PairResponse)
	b := pair[:3]
	c := pair[3:]

	r, err := http.Get(lunoPriceChartURL + "?currency=" + pair)
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(r.Body).Decode(p)

	var price float64
	for _, pair := range p.Pairs {
		if pair.BaseCode == b && pair.CounterCode == c {
			price = pair.Price
		}
	}

	return humanize.Commaf(price)
}
