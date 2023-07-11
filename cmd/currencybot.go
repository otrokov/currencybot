/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	json "encoding/json"

	telebot "gopkg.in/telebot.v3"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	Teletoken = os.Getenv("TELE_TOKEN")
)

var botCmd = &cobra.Command{
	Use:     "currencybot",
	Aliases: []string{"start"},
	Short:   "start bot",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("currencybot %s started", appVersion)

		currencybot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		currencybot.Handle(telebot.OnText, func(m telebot.Context) error {
			log.Printf(m.Message().Payload, m.Text())
			payload := m.Message().Payload
			log.Print(payload)

			switch payload {
			case "usd":
				// Monobank
				client := resty.New()
				resp, err := client.R().Get("https://api.monobank.ua/bank/currency")
				if err != nil {
					//currencybot.Reply(m, "Ошибка при получении данных с API Monobank.")
					return err
				}
				defer resp.RawResponse.Body.Close()

				var data []map[string]interface{}
				err = json.Unmarshal(resp.Body(), &data)
				/*if err != nil {
					currencybot.Reply(m, "Ошибка при обработке данных с API Monobank.")
					return err
				}*/

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
				resp, err = client.R().Get("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5")
				/*if err != nil {
					currencybot.Reply(m, "Ошибка при получении данных с API Privatbank.")
					return err
				}*/
				defer resp.RawResponse.Body.Close()

				var dataPrivat []map[string]interface{}
				err = json.Unmarshal(resp.Body(), &dataPrivat)
				/*if err != nil {
					currencybot.Reply(m, "Ошибка при обработке данных с API Privatbank.")
					return err
				}*/

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

				response := fmt.Sprintf("Курс доллара:\n\n"+"Monobank: %.2f / %.2f\n"+"Privatbank: %.2f / %.2f",
					rateBuyMono, rateSellMono, rateBuyPrivat, rateSellPrivat)
				log.Printf("The response is %s", response)
				currencybot.Send(m.Sender(), response)

			}

			return err
		})

		currencybot.Start()

	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
