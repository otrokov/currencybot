package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/tucnak/telebot.v2"
)

// botCmd represents the bot command

var (
	Teletoken = os.Getenv("TELE_TOKEN")
)

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Start the Telegram bot",

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("bot %s started", appVersion)
		fmt.Printf("TAPI= %s", Teletoken)

		bot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		/*Run: func(cmd *cobra.Command, args []string) {
		// Create a new Telegram bot instance
		bot, err := telebot.NewBot(telebot.Settings{
			Token: "TEL_TOKEN",
		})
		if err != nil {
			log.Fatal(err)
		}*/

		// Handle the /start command
		bot.Handle("/start", func(m *telebot.Message) {
			bot.Reply(m, "Привет! Отправь мне сообщение для получения валютных курсов.")
		})

		// Handle all other messages
		bot.Handle(telebot.OnText, func(m *telebot.Message) {
			// Monobank
			resp, err := requests.Get("https://api.monobank.ua/bank/currency")
			if err != nil {
				bot.Reply(m, "Ошибка при получении данных с API Monobank.")
				return
			}
			defer resp.Close()

			var data []map[string]interface{}
			err = json.NewDecoder(resp).Decode(&data)
			if err != nil {
				bot.Reply(m, "Ошибка при обработке данных с API Monobank.")
				return
			}

			var rateBuyMono, rateSellMono float64
			for _, item := range data {
				currencyCodeA := item["currencyCodeA"].(float64)
				currencyCodeB := item["currencyCodeB"].(float64)
				if currencyCodeA == 840 && currencyCodeB == 980 {
					rateBuyMono = item["rateBuy"].(float64)
					rateSellMono = item["rateSell"].(float64)
					break
				}
			}

			// Privatbank
			resp, err = requests.Get("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5")
			if err != nil {
				bot.Reply(m, "Ошибка при получении данных с API Privatbank.")
				return
			}
			defer resp.Close()

			var dataPrivat []map[string]interface{}
			err = json.NewDecoder(resp).Decode(&dataPrivat)
			if err != nil {
				bot.Reply(m, "Ошибка при обработке данных с API Privatbank.")
				return
			}

			var rateBuyPrivat, rateSellPrivat float64
			for _, item := range dataPrivat {
				ccy := item["ccy"].(string)
				baseCcy := item["base_ccy"].(string)
				if ccy == "USD" && baseCcy == "UAH" {
					rateBuyPrivat, _ = strconv.ParseFloat(item["buy"].(string), 64)
					rateSellPrivat, _ = strconv.ParseFloat(item["sale"].(string), 64)
					break
				}
			}

			response := fmt.Sprintf("Курс доллара:\n\n"+
				"Monobank: %.2f / %.2f\n"+
				"Privatbank: %.2f / %.2f",
				rateBuyMono, rateSellMono, rateBuyPrivat, rateSellPrivat)

			bot.Reply(m, response)
		})

		// Start the bot
		bot.Start()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
