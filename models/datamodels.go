package models;

import (
	// "time"

	"gorm.io/gorm"

)

type UserModel struct{
	gorm.Model
	Username string `gorm:"primarykey"`
	Password string 
	Balance float32


}


type Stock struct{
	Ticker string
	Name string
	Price float32
	MarketCap float32

}

func NewUser(uname string, pword string)* UserModel{
	return &UserModel{
		Username: uname,
		Password: pword,
	}

}