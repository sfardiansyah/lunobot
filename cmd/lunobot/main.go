package main

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/yanzay/tbot"
	"github.com/yanzay/tbot/model"
)

// PriceCandle ...
type PriceCandle struct {
	Candles []Candle `json:"candles"`
}

// Candle ...
type Candle struct {
	Timestamp time.Time `json:"timestamp"`
	High      float64   `json:"high,string"`
	Low       float64   `json:"low,string"`
}

// Ticker ...
type Ticker struct {
	Ask       float64   `json:"ask,string"`
	Timestamp time.Time `json:"timestamp"`
	Bid       float64   `json:"bid,string"`
	Volume    float64   `json:"rolling_24_hour_volume,string"`
	LastTrade float64   `json:"last_trade,string"`
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

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"), tbot.WithWebhook("https://luno-bot.herokuapp.com/", ":"+os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/update", "update")
	bot.Handle("/help", fileReader("assets/help.txt"))
	bot.HandleFunc("/infoluno", infoHandler)
	bot.HandleFunc("/start", startHandler)
	bot.HandleFunc("/fee", func(m *tbot.Message) {
		feeHandler(m, fileReader("assets/fee.txt"), "Kunjungi Rincian Biaya LUNO")
	})
	bot.HandleFunc("/convert", func(m *tbot.Message) {
		feeHandler(m, fileReader("assets/convert.txt"), "Kunjungi LUNO Price Chart")
	})

	// bot.HandleDefault(defaultHandler)

	if err := bot.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func infoHandler(m *tbot.Message) {
	btc := getPrice("XBTIDR")
	eth := getPrice("ETHIDR")
	btcH, btcL := getHiLo("XBTIDR")
	ethH, ethL := getHiLo("ETHIDR")

	m.Replyf(fileReader("assets/info.txt"), btc, btcH, btcL, eth, ethH, ethL)
}

func getHiLo(pair string) (string, string) {
	p := new(PriceCandle)
	date := time.Now().Add(-24 * time.Hour).Unix()

	log.Println(lunoChartCandleURL + "?pair=" + pair + "&since=" + strconv.FormatInt(date, 10))
	r, err := http.Get(lunoChartCandleURL + "?pair=" + pair + "&since=" + strconv.FormatInt(date, 10))
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewDecoder(r.Body).Decode(p); err != nil {
		log.Fatal(err)
	}

	var hi, lo float64
	for _, candle := range p.Candles {
		log.Println(candle.High, candle.Low)
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

func feeHandler(m *tbot.Message, t, helper string) {
	str := strings.Split(t, "||")
	btn := []map[string]string{map[string]string{helper: str[1]}}

	m.ReplyInlineKeyboard(str[0], btn, tbot.WithURLInlineButtons)
}

// func defaultHandler(m *tbot.Message) {
// 	// if len(m.Vars["new_chat_members"]) > 0 {
// 	// 	m.Reply("Halo!")
// 	// }
// 	m.Reply(m.Data)
// }

func startHandler(m *tbot.Message) {
	if m.ChatType == model.ChatTypePrivate {
		m.Reply(fileReader("assets/start.txt"))
		m.Reply(fileReader("assets/help.txt"))
	}
}

func fileReader(dir string) string {
	b, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println(err)
	}

	return html.UnescapeString(string(b))
}
