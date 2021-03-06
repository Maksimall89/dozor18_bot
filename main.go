package main

import (
	"fmt"
	"github.com/Maksimall89/dozor18_bot/src"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte("Hi there! I'm telegram bot @dozor18_bot. My owner @maksimall89"))
}

func main() {
	// init configuration
	configuration := src.Config{}
	configuration.Init()

	// monitoring
	http.HandleFunc("/", MainHandler)
	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%s", configuration.Port), nil)
	}()

	// configuration bot
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Println(err)
	}

	bot.Token = configuration.TelegramBotToken
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	defer log.Println("Bot off!.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0

	// config webhook
	updates := bot.ListenForWebhook(fmt.Sprintf("/read%s", configuration.ListenPath))
	if err != nil {
		log.Println(err)
		log.Println("Failed to get updates")
	}

	// get info from DB
	dbConfig := src.Config{}
	dbConfig.Init()

	var str string
	var command string

	// read from channel
	for update := range updates {

		if update.Message == nil {
			if update.CallbackQuery != nil {
				command = update.CallbackQuery.Data
				update.Message = update.CallbackQuery.Message
				update.Message.From = update.CallbackQuery.From
			} else {
				continue
			}
		} else {
			if update.Message.From.IsBot {
				continue
			}
			command = strings.ToLower(update.Message.Command())
		}

		if update.Message.From.ID == configuration.OwnID {
			switch command {
			case "reset", "restart":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBTruncTables(strings.TrimSpace(strings.ToLower(update.Message.CommandArguments()))), 0, bot, "admin")
				continue
			case "show":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesAll(dbConfig), 0, bot, "admin")
				continue
			case "add":
				arrCodes := strings.Split(update.Message.Text, "\n")
				for _, code := range arrCodes {
					code = strings.Replace(code, "/add", "", -1)
					str = dbConfig.DBInsertCodesRight(code)
					_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot, "admin")
				}
				continue
			case "update":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBUpdateCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot, "admin")
				continue
			case "delete":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot, "admin")
				continue
			case "listteams":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowTeams(true, dbConfig), update.Message.MessageID, bot, "admin")
				continue
			case "createdb":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBCreateTables(), 0, bot, "admin")
				continue
			case "createtask":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.CreateTask(update.Message, dbConfig), 0, bot, "admin")
				continue
			case "updatetask":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBUpdateTask(update.Message.CommandArguments()), update.Message.MessageID, bot, "admin")
				continue
			case "deletetask":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteTask(update.Message.CommandArguments()), update.Message.MessageID, bot, "admin")
				continue
			}
		}

		switch command {
		case "codes":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesMy(update.Message, dbConfig), update.Message.MessageID, bot, "main")
		case "generate", "gen":
			strArr := strings.Split(update.Message.CommandArguments(), ",")
			number, err := strconv.Atoi(strArr[0])
			if err != nil {
				_ = src.SendMessageTelegram(update.Message.Chat.ID, "&#10071;Не по формату:\n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>", update.Message.MessageID, bot, "main")
				continue
			}
			switch len(strArr) {
			case 1:
				str = src.CodeGen("", "", number, src.NameFileWords)
			case 3:
				str = src.CodeGen(strArr[1], strArr[2], number, src.NameFileWords)
			}
			_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot, "main")
		case "text", "task":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetTasks(update.Message, dbConfig), update.Message.MessageID, bot, "main")
		case "help", "start":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetListHelps(update.Message.From, configuration.OwnID), update.Message.MessageID, bot, "main")
		case "create":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.CreateTeam(update.Message, dbConfig), update.Message.MessageID, bot, "team")
		case "listusers":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowUsers(update.Message, false, dbConfig), update.Message.MessageID, bot, "team")
		case "list":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowUsers(update.Message, true, dbConfig), update.Message.MessageID, bot, "team")
		case "join":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.JoinTeam(update.Message, dbConfig), update.Message.MessageID, bot, "team")
		case "leave":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteUser(update.Message.From.ID), update.Message.MessageID, bot, "team")
		case "invite":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetInvite(update.Message, dbConfig), update.Message.MessageID, bot, "team")
		case "teams":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowTeams(false, dbConfig), update.Message.MessageID, bot, "team")
		default:
			if strings.HasPrefix(update.Message.Text, "/") {
				_ = src.SendMessageTelegram(update.Message.Chat.ID, "&#9940;I don't know what you want. But you can use /help", update.Message.MessageID, bot, "main")
				continue
			}
			src.CheckCode(update.Message, bot, dbConfig)
		}
	}
}
