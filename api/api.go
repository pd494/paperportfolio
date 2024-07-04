package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"stock-portfolio-cli/database"
	"stock-portfolio-cli/models"
	"sync"
	"time"

	"strings"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"errors"
)

var mu sync.Mutex
var Red = "\033[31m" 
var Blue = "\033[34m" 
var Reset = "\033[0m" 
var Cyan = "\033[36m" 
var Green = "\033[32m" 
var White = "\033[97m"




func GetCurrentPrice(ticker string) (finnhub.Quote, error){

	err:= godotenv.Load("/Users/prasanthdendukuri/Projects/stockportfolio/.env")
	if err!= nil{
		fmt.Println("couldnt load env")
	}

	// secretkey:= os.Getenv("SECRET_KEY")
	apikey:= os.Getenv("API_KEY")

	cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apikey)
    finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi
	res, _, err := finnhubClient.Quote(context.Background()).Symbol(ticker).Execute()
	if err != nil{
		fmt.Println("ticker api call execution failed")
		return finnhub.Quote{}, errors.New("ticker not found")

	}
	
	return res, nil;


}


func NetGain(user *models.UserModel) {
	time.Sleep(15* time.Second)
	apikey := os.Getenv("API_KEY")
	w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token=" + apikey, nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	stocks, _ := database.GetStocksByUsername(user.Username)
	var symbols []string
	for _, stock := range *stocks {
		symbols = append(symbols, stock.Ticker)
	}

	for i, symbol := range symbols {
		symbols[i] = strings.ToUpper(symbol)
	}
	symbols = append(symbols, "IC MARKETS:1")

	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		err := w.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Failed to subscribe to symbol:", s, err)
		}
	}

	type Message struct {
		Data []struct {
			P float64 `json:"p"`
			S string  `json:"s"`
			T int64   `json:"t"`
			V float64 `json:"v"`
		} `json:"data"`
		Type string `json:"type"`
	}

	var msg Message
	for {

		time.Sleep(15 * time.Second)
		err := w.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON message from WebSocket:", err)
			return
		}
		if msg.Type == "trade" {
			for _, trade := range msg.Data {
				for i := range *stocks {
					if (*stocks)[i].Ticker == trade.S {
						mu.Lock()
						(*stocks)[i].Gain = float32((trade.P - (*stocks)[i].AveragePrice) * (*stocks)[i].QuantityOwned)

						if err := database.GetDB().Save(&(*stocks)[i]).Error; err != nil {
							fmt.Println(Red, "Failed to save stock gain:", err)
						} else {
							fmt.Println(Green, "Updated gain for stock:", trade.S, "New Gain:", (*stocks)[i].Gain)
						}

						mu.Unlock()
					}
				}
			}
		}
	}
}











