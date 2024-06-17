package database

import (
	_ "database/sql"
	"fmt"
	"os"
	"stock-portfolio-cli/models" // Adjust the import path according to your project structure
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var db* gorm.DB

func DBInit(){

	err:= godotenv.Load("/Users/prasanthdendukuri/Projects/stockportfolio/.env")

	if err!= nil{
		fmt.Println("Error loading .env file:", err)

	}

		host:= os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user:= os.Getenv("DB_USER")
		password:= os.Getenv("DB_PASS")
		dbname:= os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    
		
		dbConn, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		db = dbConn
		

		fmt.Println("connection to db succesful")


}

func ValidateLogin(username string, password string) (bool, error) {
	var user models.UserModel
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil{
		return false, result.Error
	}
	if user.Password != password{
		return false, fmt.Errorf("incorrect password")
	}

	return true, nil
}