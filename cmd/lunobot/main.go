package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sfardiansyah/tbot"
	"github.com/sfardiansyah/tbot/model"
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

func main() {
	app, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	app.Debug = true

	info, err := app.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsSet() {
		_, err = app.SetWebhook(tgbotapi.NewWebhook("https://luno-bot.herokuapp.com:443"))
		if err != nil {
			log.Fatal(err)
		}
	}

	// h := rest.Handler(app)
	updates := app.ListenForWebhook("/")
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		log.Printf("%+v\n", update.Message)
	}

	// bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"), tbot.WithWebhook("https://luno-bot.herokuapp.com/", ":"+os.Getenv("PORT")))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// bot.Handle("/update", "update")
	// bot.Handle("/help", fileReader("assets/help.txt"))
	// bot.HandleFunc("/infoluno", func(m *tbot.Message) {
	// 	replyWithInline(m, getInfo(), "Buka LUNO Wallet")
	// })
	// bot.HandleFunc("/start", startHandler)
	// bot.HandleFunc("/fee", func(m *tbot.Message) {
	// 	replyWithInline(m, fileReader("assets/fee.txt"), "Kunjungi Rincian Biaya LUNO")
	// })
	// bot.HandleFunc("/convert", func(m *tbot.Message) {
	// 	replyWithInline(m, fileReader("assets/convert.txt"), "Kunjungi LUNO Price Chart")
	// })

	// bot.HandleDefault(defaultHandler)

	// if err := bot.ListenAndServe(); err != nil {
	// 	log.Fatal(err)
	// }
}

func getInfo() string {
	idTime, _ := time.LoadLocation("Asia/Jakarta")
	date := time.Now().In(idTime).Format(time.RFC1123)

	btc := getPrice("XBTIDR")
	eth := getPrice("ETHIDR")
	btcH, btcL := getHiLo("XBTIDR")
	ethH, ethL := getHiLo("ETHIDR")
	btcVol := getVolume("XBTIDR")
	ethVol := getVolume("ETHXBT")

	return fmt.Sprintf(fileReader("assets/info.txt"), date, btc, btcVol, btcH, btcL, eth, ethVol, ethH, ethL)
}

func getVolume(pair string) string {
	t := new(Ticker)

	r, err := http.Get(lunoAPIURL + "?pair=" + pair)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Fatal(err)
	}

	return humanize.CommafWithDigits(t.Volume, 8)
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

func replyWithInline(m *tbot.Message, t, helper string) {
	str := strings.Split(t, "||")
	btn := []map[string]string{map[string]string{helper: str[1]}}

	m.ReplyInlineKeyboard(str[0], btn, tbot.WithURLInlineButtons)
}

func defaultHandler(m *tbot.Message) {
	m.Reply(m.Text())
	log.Println(m.ChatID, m.ChatType)
}

func startHandler(m *tbot.Message) {
	if m.ChatType == model.ChatTypePrivate {
		m.Reply(fileReader("assets/start.txt"), tbot.WithMarkdown)
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
