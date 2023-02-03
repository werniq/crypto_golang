package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tuxofil/golang-kunaio/src/kunaio"
	"github.com/yanzay/tbot/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Market              string
	PriceBID            int
	OrderBookVolume     int
	PriceASK            int
	VolumeASK           float64
	PriceChangeCurrency int
	PriceChangePercent  float32
	LatestPrice         int
	DayVolume           int
	DayMaximumPrice     int
	DayMinimumPrice     int
}

type application struct {
	client *tbot.Client
}

var (
	errorLog  *log.Logger
	app       application
	bot       *tbot.Server
	token     string
	publicKey = "qfy9hSU5t2Zck41rp7FwGbm3iXAeUft9AvDqqBgm"
	secretKey = "h2xV8CVUa96qxei1pFCxmJcm08Si5tHbH2Ir5va9"
	tgBotKey  = "6145823087:AAHg69rPpNI1y1fK90kDlkCeILO0pDO1Ta0"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

func main() {
	bot := tbot.New(tgBotKey)

	app.client = bot.Client()
	//bot.HandleMessage("/stats", app.)
	//markets := kunaio.SupportedMarkets()
	//fmt.Println(markets)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Lshortfile|log.Ltime|log.Lmsgprefix)
	bot.HandleMessage("/start", app.Start)
	bot.HandleMessage("/info", app.Info)
	bot.HandleMessage("/myinfo", app.MyInfo)
	bot.HandleMessage("/price", app.Price)
	bot.HandleMessage("/neworder", app.Buy)
	bot.HandleMessage("/calc", app.CalculateProfit)
	//fmt.Println(bot)
	errorLog.Println(bot.Start())
	fmt.Println("Discord bot is currently running. Press CTRL-C to stop.")

}

func (app *application) MyInfo(m *tbot.Message) {
	userInfo, err := kunaio.GetUserInfo(publicKey, secretKey)
	if err != nil {
		errorLog.Println(err)
	}
	for i, account := range userInfo.Accounts {
		if userInfo.Accounts[i].Balance != float64(0) {
			msg := fmt.Sprintf("Currency: %s; Balance: %f; Locked: %f\n",
				account.Currency, account.Balance, account.Locked)
			app.client.SendMessage(m.Chat.ID, msg)
		}
	}
}

func (app *application) Start(m *tbot.Message) {
	msg := "Hey! I am bot for trading crypto currency in [Kuna] exchange"
	app.client.SendMessage(m.Chat.ID, msg)
}

func (app *application) CalculateProfit(m *tbot.Message) {
	args := strings.Split(m.Text, " ")
	command := args[0]
	//fmt.Println(args)

	if command != "" {
		if len(args) == 4 {
			fmt.Println(args)
			first, err := strconv.Atoi(args[1])
			if err != nil {
				errorLog.Println(err)
			}
			second, err := strconv.Atoi(args[2])
			if err != nil {
				errorLog.Println(err)
			}
			var gainPercentage float32
			gainPercentage = float32(first / second)
			fmt.Println(gainPercentage)
			gainAmount := gainPercentage
			if gainPercentage > 0 {
				app.client.SendMessage(m.Chat.ID, fmt.Sprintf("Approximate gain percentage is: %f", gainPercentage))
				app.client.SendMessage(m.Chat.ID, fmt.Sprintf("Approximate gain amount is: %f", gainAmount))
			} else {
				app.client.SendMessage(m.Chat.ID, fmt.Sprintf("You will lose %f", gainPercentage))
			}
		} else {
			app.client.SendMessage(m.Chat.ID, fmt.Sprintf("You should provide following arguments [buyPrice] [sellPrice] [Amount of $]"))
		}
	} else {
		app.client.SendMessage(m.Chat.ID, "Error 1 ")
	}

}

func (app *application) Price(m *tbot.Message) {
	currentMarket := "ethuah"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.kuna.io/v3/tickers?symbols=ethuah", nil)
	req.Header.Set("Accept", "application/json")

	res, _ := client.Do(req)

	if err != nil {
		errorLog.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	app.client.SendMessage(m.Chat.ID, string(body))
	//fmt.Println(response)
	if err != nil {
		errorLog.Println(err)
		return
	}

	if err != nil {
		app.client.SendMessage(m.Chat.ID, fmt.Sprintf("Error getting stats for [%s] market, %v", strings.ToUpper(currentMarket), err))
		return
	}
}

func (app *application) Info(m *tbot.Message) {
	market := "ethuah"
	orderBook, err := kunaio.GetOrderBook(market)
	if err != nil {
		log.Fatalf("Error getting order book: %v", err)
		return
	}
	for _, order := range orderBook.Bids {
		msg := fmt.Sprintf(`
			Market: %s 
	 		Price: %f
			Avarage price: %f 
			OrderType: %s 
			Trades count: %d
			Volume %f
			`, order.Market, order.Price, order.AvgPrice, order.OrdType, order.TradesCount, order.Volume)
		app.client.SendMessage(m.Chat.ID, msg)
	}
	for _, order := range orderBook.Asks {
		msg := fmt.Sprintf(`
			Market: %s 
	 		Price: %f
			Avarage price: %f 
			OrderType: %s 
			Trades count: %d
			Volume %f
			`, order.Market, order.Price, order.AvgPrice, order.OrdType, order.TradesCount, order.Volume)
		app.client.SendMessage(m.Chat.ID, msg)
	}
}

func (app *application) Buy(m *tbot.Message) {
	market := "ETHUAH"
	stats, err := kunaio.GetLatestStats(market)
	if err != nil {
		log.Fatalf("Error getting ETHUAH stats: %v", err)
		return
	}

	bidTime := time.Now()
	msg := fmt.Sprintf("Bid time is %T. Current ETH Price is: %f", bidTime, stats.Buy)
	app.client.SendMessage(m.Chat.ID, msg)

	//side := "buy"
	userStats, err := kunaio.GetUserInfo(publicKey, secretKey)
	//balance := int(userStats.Accounts[3].Balance / 20)

	fmt.Println(userStats.Accounts[3])

	//order, err := kunaio.NewOrder(publicKey, secretKey, market, side, float64(balance), stats.Buy)
	//fmt.Printf("Order created: %+v\n", order)
}
