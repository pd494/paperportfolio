package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/howeyc/gopass"
	"stock-portfolio-cli/database"
	"stock-portfolio-cli/models"


)

var loggedIn bool = false
var currUser *models.UserModel; 

var rootCmd = &cobra.Command{
	Use:   "stock-cli",
	Short: "A CLI application for stock trading",
}

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login to your account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("Enter password for %s: ", username)
		passwd, _ := gopass.GetPasswdMasked()

		success, err := database.ValidateLogin(username, string(passwd))
		if err != nil {
			fmt.Println("Login error:", err)
			return
		}
		if success {
			loggedIn = true
			fmt.Println("Login successful!")
		} else {
			fmt.Println("Login failed: Incorrect username or password")
		}
	},
}

var buyCmd = &cobra.Command{
	Use:   "buy [stock]",
	Short: "Buy a stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !loggedIn {
			fmt.Println("You need to be logged in to buy stocks.")
			return
		}
		stock := args[0]
		fmt.Printf("Buying stock: %s\n", stock)
		// Implement logic to buy the specified stock
	},
}

var sellCmd = &cobra.Command{
	Use:   "sell [stock]",
	Short: "Sell a stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if (!loggedIn) {
			fmt.Println("You need to be logged in to sell stocks.")
			return
		}
		stock := args[0]
		fmt.Printf("Selling stock: %s\n", stock)
		// Implement logic to sell the specified stock
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
		fmt.Println("Viewing portfolio:")
		// Implement logic to view the user's portfolio
	},
}

func main() {
	database.DBInit()

	// Register commands
	rootCmd.AddCommand(loginCmd, buyCmd, sellCmd, portfolioCmd)

	go func() {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	// Start continuous input loop
	inputLoop()
}

func inputLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
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
