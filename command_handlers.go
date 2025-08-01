package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MenuCommand = "Menu"
const SubscribeCommand = "Subscribe"
const UnsubcribeCommand = "Unsubscribe"
const GoonCommand = "Goon"
const DownloadCommand = "Download"
const StatusCommand = "Status"

func (a *Application) HandleMenu(chatId int64) error {
	text := "Please Select your option 👇"
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = tgbotapi.ModeHTML
	// Keyboard layout for the second menu. Two buttons, one per row
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Subscribe 🚀", SubscribeCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Unsubscribe 😩", UnsubcribeCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Goon 🍆💦", GoonCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Download ⬇️🌽", DownloadCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Status 🔎", StatusCommand),
		),
	)
	msg.ReplyMarkup = markup
	_, err := a.bot.Send(msg)
	return err
}

func (a *Application) HandleGoon(chatId int64) error {
	err := a.sendText(chatId, "Goon material incoming! Prepare your dicks! 🍆💦😩")

	if err != nil {
		return err
	}

	err = a.sendRedditPorn(chatId)

	if err != nil {
		err := a.sendText(chatId, "There were some issue during the scraping process. Please notify the dev to fix it.")
		return err
	}

	err = a.sendText(chatId, "That's it for now hope you like it. 😉")

	return err
}

func (a *Application) HandleStatus(chatId int64) error {
	statusText := "🤖 Bot Status: Online\n📊 System: Operational\n🔄 Last Update: Active"
	msg := tgbotapi.NewMessage(chatId, statusText)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := a.bot.Send(msg)
	return err
}

func (a *Application) HandleDownload(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, "📥 Download feature is under development. Coming soon!")
	_, err := a.bot.Send(msg)
	return err
}

func (a *Application) HandleSubscribe(chatId int64) error {
	// TODO: Implement subscription logic (save user to database, etc.)
	msg := tgbotapi.NewMessage(chatId, "🚀 Successfully subscribed! You'll receive updates automatically.")
	_, err := a.bot.Send(msg)
	return err
}

func (a *Application) HandleUnsubscribe(chatId int64) error {
	// TODO: Implement unsubscription logic (remove user from database, etc.)
	msg := tgbotapi.NewMessage(chatId, "😩 Successfully unsubscribed. You won't receive updates anymore.")
	_, err := a.bot.Send(msg)
	return err
}
