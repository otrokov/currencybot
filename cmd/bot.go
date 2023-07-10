import telebot
import requests
import json

bot = telebot.TeleBot('TELETOKEN')

@bot.message_handler(commands=['start'])
def handle_start(message):
    bot.reply_to(message, 'Привет! Отправь мне сообщение для получения валютных курсов.')

@bot.message_handler(func=lambda message: True)
def handle_message(message):
    try:
        # Monobank
        response = requests.get('https://api.monobank.ua/bank/currency')
        json_data = response.json()

        usd_rateSellMono = next(item for item in json_data if item['currencyCodeA'] == 840 and item['currencyCodeB'] == 980)
        rateSellMono = usd_rateSellMono['rateSell']
        usd_rateBuyMono = next(item for item in json_data if item['currencyCodeA'] == 840 and item['currencyCodeB'] == 980)
        rateBuyMono = usd_rateBuyMono['rateBuy']
        text_mono = f'Курс доллара Mono: {rateBuyMono} / {rateSellMono}'

        # Privatbank
        response = requests.get('https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5')
        json_data = response.json()

        usd_rateSellPrivat = next(item for item in json_data if item['ccy'] == 'USD' and item['base_ccy'] == 'UAH')
        rateSellPrivat = usd_rateSellPrivat['sale']
        usd_rateBuyPrivat = next(item for item in json_data if item['ccy'] == 'USD' and item['base_ccy'] == 'UAH')
        rateBuyPrivat = usd_rateBuyPrivat['buy']
        text_privat = f'Курс доллара Privat: {rateBuyPrivat} / {rateSellPrivat}'

        bot.reply_to(message, f'{text_mono}\n{text_privat}')
    except requests.exceptions.RequestException:
        bot.reply_to(message, 'Ошибка при получении данных с API.')

bot.polling()

