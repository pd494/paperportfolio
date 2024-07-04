package database

import (
	_ "database/sql"
	"fmt"
	"os"
	"stock-portfolio-cli/models" // Adjust the import path according to your project structure

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "reflect"
)

var db *gorm.DB
var Red = "\033[31m"
var Blue = "\033[34m"
var Reset = "\033[0m"
var Cyan = "\033[36m"
var Green = "\033[32m"

func GetDB() *gorm.DB{
	return db
}
func DBInit() {

	err := godotenv.Load("/Users/prasanthdendukuri/Projects/stockportfolio/.env")

	if err != nil {
		fmt.Println("Error loading .env file:", err)

	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	if err != nil {
		panic("failed to perform migrations: " + err.Error())
	}

	dbConn, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db = dbConn
	err = db.AutoMigrate(&models.UserModel{}, &models.Stock{})
	db.Migrator().DropColumn(&models.Stock{}, "name")

	fmt.Println("connection to db succesful")

}

func ValidateLogin(username string, password string) (*models.UserModel, error) {
	var user models.UserModel
	result := db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	if user.Password != password {
		return nil, fmt.Errorf(Red + "incorrect password" + Reset)
	}

	return &user, nil
}

func AlreadyExists(username string) bool {
	var count int64
	result := db.Model(&models.UserModel{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false
	}
	return count > 0
}

func RegisterUser(username string, password string, balance float64) *models.UserModel {
	NewUser := models.UserModel{
		Username: username,
		Password: password,
		Balance:  balance,
	}

	result := db.Create(&NewUser)
	if result.Error != nil {
		panic("failed to create user: " + result.Error.Error())
	}

	return &NewUser

}

func GetStocksByUsername(username string) (*[]models.Stock, error) {
	var user models.UserModel
	result := db.Preload("Stocks").Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user.Stocks, nil
}

func BuyStock(quantity float64, ticker string, user *models.UserModel, res finnhub.Quote) {
	price := float64(*res.C)
	if quantity*price > user.Balance {
		fmt.Println(Red+"Not enough balance to buy "+ticker, Reset)
		return
	} else {

		stocks, _ := GetStocksByUsername(user.Username)

		for i := range *stocks {
			if (*stocks)[i].Ticker == ticker {
				(*stocks)[i].AveragePrice = (((*stocks)[i].AveragePrice * (*stocks)[i].QuantityOwned) + (price * quantity)) / ((*stocks)[i].QuantityOwned + quantity)

				(*stocks)[i].QuantityOwned += quantity
				(*stocks)[i].TotalWorth += quantity * price

				user.Balance -= quantity * price

				if err := db.Save(user).Error; err != nil {
					fmt.Println(Red, "Failed to save user balance:", err)
					return
				}

				if err := db.Save(&(*stocks)[i]).Error; err != nil {
					fmt.Println(Red, "Failed to save stock:", err)
					return
				}

				fmt.Println(Green, "Transaction successful!", Blue)
				fmt.Println("Balance remaining:", user.Balance)
				fmt.Print(Reset)
				return
			}
		}



		newStock := models.Stock{
			Ticker:        ticker,
			TotalWorth:    float64(price) * quantity,
			AveragePrice:  float64(price),
			QuantityOwned: quantity,
			Gain:          0,
		}
		user.Balance -= quantity * price

		user.Stocks = append(user.Stocks, newStock)
		savestock := db.Save(user)
		if savestock.Error != nil {
			fmt.Println(Red, "unable to save changes at this time")
			return
		}
		fmt.Println(Green, "Transaction succesful!")
		fmt.Println("Balance remaining:", user.Balance)
		fmt.Print(Blue)

	}

}

func PrintPortfolio(username string) {
	var user models.UserModel
	db.Where("username = ?", username).Preload("Stocks").First(&user)

	fmt.Println(Cyan, "Portfolio:")
	for _, stock := range user.Stocks {
		fmt.Printf("Ticker: %s, Quantity: %.2f, Total Worth: %.2f, Average Price: %.2f, Gain: %.2f\n", stock.Ticker, stock.QuantityOwned, stock.TotalWorth, stock.AveragePrice, stock.Gain)
	}
	fmt.Print(Reset)
}


func SellStock(quantity float64, ticker string, user *models.UserModel, res finnhub.Quote) {
	price := float64(*res.C)

	if err := db.Where("username = ?", user.Username).Preload("Stocks").First(&user).Error; err != nil {
		fmt.Println(Red + "Failed to fetch user's stocks:", err)
		return
	}

	var stock models.Stock
	found := false
	for i := range user.Stocks {
		if user.Stocks[i].Ticker == ticker {
			stock = user.Stocks[i]
			found = true
			break
		}
	}

	if !found {
		fmt.Println(Red + "Stock not found in portfolio:" + ticker)
		return
	}

	if stock.QuantityOwned < quantity {
		fmt.Println(Red + "Not enough quantity to sell:" + ticker)
		return
	}

	gain := (quantity * price) - (quantity * stock.AveragePrice)
	stock.Gain += float32(gain)

	stock.QuantityOwned -= quantity
	stock.TotalWorth -= quantity * price
	user.Balance += quantity * price

	if err := db.Save(user).Error; err != nil {
		fmt.Println(Red, "Failed to save user balance:", err)
		return
	}

	if err := db.Save(&stock).Error; err != nil {
		fmt.Println(Red, "Failed to save stock:", err)
		return
	}

	fmt.Println(Green, "Transaction successful!", Blue)
	fmt.Println("Balance remaining:", user.Balance)
	fmt.Print(Reset)
}