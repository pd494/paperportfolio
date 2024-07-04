package main

import (
	"bufio"
	"fmt"
	"os"
	_ "os/user"
	"strings"

	"stock-portfolio-cli/api"
	"stock-portfolio-cli/database"
	"stock-portfolio-cli/models"
	"strconv"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var loggedIn bool = false
var currUser *models.UserModel; 
var Red = "\033[31m" 
var Blue = "\033[34m" 
var Reset = "\033[0m" 
var Cyan = "\033[36m" 
var Green = "\033[32m" 
var White = "\033[97m"


var rootCmd = &cobra.Command{
	// Use:   "stock-cli",
	Short: "A CLI application for stock trading",
	SilenceUsage: true, // This will prevent the default usage message from being printed

}

var infoCmd = &cobra.Command{
	Use: "info [ticker]",
	Short: "get info about stock tickers. provide one ticker", 
	Args: cobra.ExactArgs(1), 
	Run: func(cmd *cobra.Command, args [] string){
		result, err := api.GetCurrentPrice(args[0])
		
		if err != nil{
			fmt.Print(Blue + fmt.Sprintf("ticker not found, try again"))

		}else if *result.C == 0.0 && *result.H == 0.0 && *result.L == 0.0 && *result.O == 0.0 {
			fmt.Println(Blue + "Invalid ticker data received, please try again" + Reset)
		}else
		{
			fmt.Print(Blue + fmt.Sprintf("%v", "Current Price: "))
			fmt.Println(Blue + fmt.Sprintf("%v", *result.C))
			fmt.Print(Blue + fmt.Sprintf("%v", "Day High: "))
			fmt.Println(Blue + fmt.Sprintf("%v", *result.H))

			fmt.Print(Blue + fmt.Sprintf("%v", "Day Low: "))
			fmt.Println(Blue + fmt.Sprintf("%v", *result.L))
			fmt.Print(Blue + fmt.Sprintf("%v", "Opening Price: "))
			fmt.Println(Blue + fmt.Sprintf("%v", *result.O))
			fmt.Print(Reset);

		}
		
	},
}

var registerCmd = &cobra.Command{
	Use:   "register [username]",
	Short: "register for a new account",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		reader := bufio.NewReader(os.Stdin)

		for database.AlreadyExists(username) {
			fmt.Println(Red + "Username already exists, try again")
			fmt.Print(Reset)
			fmt.Print("Enter a new username: ")

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(Red + "Error reading input, please try again")
				fmt.Print(Reset)
				continue
			}
			username = strings.TrimSpace(input)
		}

		fmt.Println(Cyan + "> " + "Enter password for account")
		passwd, _ := gopass.GetPasswdMasked()
		var investment float64

		fmt.Print(Cyan + "> " + "Enter how much money you want to invest: ")


		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(Red + "Error reading input, please enter a valid number")
				fmt.Print(Cyan + "> " + "Enter how much money you want to invest: ")
				continue
			}

			input = strings.TrimSpace(input)
			investment, err = strconv.ParseFloat(input, 64)
			if err == nil {
				break
			}

			fmt.Println(Red + "Error reading input, please enter a valid number")
			fmt.Print(Cyan + "> " + "Enter how much money you want to invest: ")
		}

		newUser:= database.RegisterUser(username, string(passwd), investment)
		if newUser != nil{
			currUser = newUser
			fmt.Println(Green + "Success! New account has been created" + Reset);

		}
	},
}
var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login to your account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Print(White)
		fmt.Printf("Enter password for %s: ", username)
		passwd, _ := gopass.GetPasswdMasked()

		user, err := database.ValidateLogin(username, string(passwd))
		if err != nil {
			fmt.Println("Login error:", err)
			return
		}
		if user != nil {
			loggedIn = true
			currUser = user
			fmt.Println(Green + "Login successful!")
			fmt.Println(Green,"Username:",currUser.Username,"\n", "Balance:",currUser.Balance)
			fmt.Print(Reset)
			go api.NetGain(currUser)

		} else {
			fmt.Println(Red + "Login failed: Incorrect username or password" + Reset)
		}
	},
}
var buyCmd = &cobra.Command{
	Use:   "buy [stock] [amount]",
	Short: "Buy a stock",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !loggedIn {
			fmt.Println("You need to be logged in to buy stocks.")
			return
		}
		stock := args[0]
		

		if _, err := strconv.ParseFloat(args[1], 64); err != nil {
			fmt.Println(Red + "Quantity must be a valid number",Reset)
			return
		}
		quantity, _ := strconv.ParseFloat(args[1], 64)
		result, _ := api.GetCurrentPrice(stock)
		if *result.C == 0.0 && *result.H == 0.0 && *result.L == 0.0 && *result.O == 0.0 {
			fmt.Println(Red + "Invalid ticker data received, please try again" + Reset)
			return
		}

		database.BuyStock(quantity, stock, currUser, result)

	},
}

var sellCmd = &cobra.Command{
	Use:   "sell [stock] [quantity]",
	Short: "Sell a stock",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !loggedIn {
			fmt.Println("You need to be logged in to sell stocks.")
			return
		}

		stock := args[0]
		quantity, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println(Red + "Quantity must be a valid number")
			return
		}

		result, err := api.GetCurrentPrice(stock)
		if err != nil {
			fmt.Println(Red + "Failed to get current price for " + stock + Reset)
			return
		}

		database.SellStock(quantity, stock, currUser, result)
	},
}

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "View portfolio",
	Run: func(cmd *cobra.Command, args []string) {
		if (!loggedIn) {
			fmt.Println("You need to be logged in to view your portfolio.")
			return
		}
		database.PrintPortfolio(currUser.Username)
	},
}

func main() {
	database.DBInit()

	// Register commands
	rootCmd.AddCommand(loginCmd, buyCmd, sellCmd, portfolioCmd, infoCmd, registerCmd)




	go func() {
		if err := rootCmd.Execute(); err != nil {
			fmt.Print("> ")

			fmt.Println(err)
			os.Exit(1)
		}
	}()


	// Start continuous input loop
	fmt.Println("> ")
	inputLoop()
}

func inputLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ") // Print the ">" prompt
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		if args[0] == "quit" {
			fmt.Println("Exiting...")
			return
		}

		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println("Error executing command:", err)
		}
	}
}



