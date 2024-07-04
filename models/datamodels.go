package models;

import (
	// "time"

	"gorm.io/gorm"

)

type UserModel struct{
	gorm.Model
	Username string `gorm:"primarykey"`
	Password string 
	Stocks []Stock `gorm:"foreignKey:UserID"`
	Balance float64  
	NetGain float64 

}

type Stock struct{
	gorm.Model
	UserID        uint    // Foreign key
	Ticker string 
	TotalWorth float64
	AveragePrice float64
	QuantityOwned float64
	Gain float32


}


