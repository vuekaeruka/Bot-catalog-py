package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	host     = "YOUR_HOSR"
	port     = "YOUR_PORT"
	user     = "YOUR_USER"
	password = "YOUR_PASSWORD"
	dbname   = "YOUR_DB"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}

	log.Println("Успешное подключение к базе данных!")

	bot, err := tgbotapi.NewBotAPI("YOUR_TOKEN")
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v\n", err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				startMsg := "Здравствуйте! Я YOUR_BOT. Хотите просмотреть каталог товаров?"
				button := tgbotapi.NewInlineKeyboardButtonData("Начать", "start_actions")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMsg)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)

			case "end":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо за использование бота! До свидания!"))

			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}
		}

		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data

			if callbackData == "start_actions" {
				actionMsg := "Выберите действие:"
				catalogBtn := tgbotapi.NewInlineKeyboardButtonData("Посмотреть каталог", "catalog")
				endBtn := tgbotapi.NewInlineKeyboardButtonData("Завершить работу", "end_bot")
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(catalogBtn),
					tgbotapi.NewInlineKeyboardRow(endBtn),
				)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, actionMsg)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}
			if callbackData == "catalog" {
				rows, err := db.Query("SELECT id, name FROM products")
				if err != nil {
					log.Printf("Ошибка при запросе к базе данных: %v\n", err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при получении данных."))
					continue
				}
				defer rows.Close()

				var buttons [][]tgbotapi.InlineKeyboardButton
				for rows.Next() {
					var id int
					var name string
					if err := rows.Scan(&id, &name); err != nil {
						log.Printf("Ошибка при считывании товара: %v\n", err)
						continue
					}
					button := tgbotapi.NewInlineKeyboardButtonData(name, fmt.Sprintf("view:%d", id))
					row := tgbotapi.NewInlineKeyboardRow(button)
					buttons = append(buttons, row)
				}

				keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите товар:")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}

			if len(callbackData) > 0 && callbackData[:5] == "view:" {
				id := callbackData[5:]
				var name string
				var price int
				var imageUrls []string

				err := db.QueryRow("SELECT name, price FROM products WHERE id = $1", id).Scan(&name, &price)
				if err != nil {
					log.Printf("Ошибка при запросе к базе данных: %v\n", err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при получении данных."))
					continue
				}

				rows, err := db.Query("SELECT image_url FROM product_images WHERE product_id = $1", id)
				if err != nil {
					log.Printf("Ошибка при запросе изображений: %v\n", err)
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка при получении изображений."))
					continue
				}
				defer rows.Close()

				for rows.Next() {
					var imageUrl string
					if err := rows.Scan(&imageUrl); err != nil {
						log.Printf("Ошибка при считывании изображения: %v\n", err)
						continue
					}
					imageUrls = append(imageUrls, imageUrl)
				}

				msgText := fmt.Sprintf("Товар: %s\nЦена: %d рублей", name, price)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, msgText)
				bot.Send(msg)

				for _, imageUrl := range imageUrls {
					photo := tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileURL(imageUrl))
					bot.Send(photo)
				}

				managerButton := tgbotapi.NewInlineKeyboardButtonURL("Связаться с менеджером", "YOUR_LINK")
				managerKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(managerButton))

				managerMsg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы можете связаться с менеджером для покупки!:")
				managerMsg.ReplyMarkup = managerKeyboard
				bot.Send(managerMsg)

				chooseMoreBtn := tgbotapi.NewInlineKeyboardButtonData("Выбрать другой товар", "catalog")
				endBtn := tgbotapi.NewInlineKeyboardButtonData("Завершить работу", "end_bot")
				actionsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(chooseMoreBtn),
					tgbotapi.NewInlineKeyboardRow(endBtn),
				)

				actionsMsg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Хотите выбрать другой товар?")
				actionsMsg.ReplyMarkup = actionsKeyboard
				bot.Send(actionsMsg)
			}

			if callbackData == "end_bot" {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Спасибо что выбрали нас. До свидания!"))
			}
		}
	}
}