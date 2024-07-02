package api;

import (
	_"encoding/json"
	"fmt"
	_"github.com/gorilla/websocket"
	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"os"
	"github.com/joho/godotenv"
	"context"
	"errors"

)

func main() {
	GetCurrentPrice("AAPL")
	// w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token=cpp2tj1r01qn2da28f20cpp2tj1r01qn2da28f2g", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// defer w.Close()

	// symbols := []string{"AAPL", "IC MARKETS:1"}
	// for _, s := range symbols {
	// 	msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
	// 	w.WriteMessage(websocket.TextMessage, msg)
	// }

	// var msg interface{}
	// for{
	// 	err := w.ReadJSON(&msg);
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Message from server ", msg, "\n")
	// }
}

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
