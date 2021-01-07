package main

import (
	"fmt"
	"github.com/Maksimall89/dozor18_bot/src"
	"gopkg.in/telegram-bot-api.v4"
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
	// read from channel
	for update := range updates {
		if update.Message == nil || update.Message.From.IsBot {
			continue
		}

		if update.Message.From.ID == configuration.OwnID {
			switch strings.ToLower(update.Message.Command()) {
			case "reset", "restart":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBTruncTables("codes"), 0, bot)
			case "show":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesAll(dbConfig), 0, bot)
			case "add":
				arrCodes := strings.Split(update.Message.Text, "\n")
				for _, code := range arrCodes {
					code = strings.Replace(code, "/add", "", -1)
					str = dbConfig.DBInsertCodesRight(code)
					_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				}
			case "update":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBUpdateCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot)
			case "delete":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot)
			case "resetteams":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBTruncTables("teams"), 0, bot)
			case "listteams":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowTeams(true, dbConfig), update.Message.MessageID, bot)
			case "createdb":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBCreateTables(), 0, bot)
			}
		}

		switch strings.ToLower(update.Message.Command()) {
		case "codes":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesMy(update.Message, dbConfig), update.Message.MessageID, bot)
		case "generate", "gen":
			strArr := strings.Split(update.Message.CommandArguments(), ",")
			number, err := strconv.Atoi(strArr[0])
			if err != nil {
				str = "&#10071;Не по формату:\n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>"
			}
			switch len(strArr) {
			case 1:
				str = src.CodeGen("", "", number, src.NameFileWords)
			case 3:
				str = src.CodeGen(strArr[1], strArr[2], number, src.NameFileWords)
			}
			_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
		case "text":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, `Текст приквела доступен на нашем сайте <a href="http://dozor18.ru">dozor18.ru</a>.`, update.Message.MessageID, bot)
		case "help", "start":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetListHelps(update.Message.From, configuration.OwnID), update.Message.MessageID, bot)
		case "create":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.CreateTeam(update.Message, dbConfig), update.Message.MessageID, bot)
		case "listusers":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowUsers(update.Message, false, dbConfig), update.Message.MessageID, bot)
		case "list":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowUsers(update.Message, true, dbConfig), update.Message.MessageID, bot)
		case "join":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.JoinTeam(update.Message, dbConfig), update.Message.MessageID, bot)
		case "leave":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteUser(update.Message.From.ID), update.Message.MessageID, bot)
		case "invite":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetInvite(update.Message, dbConfig), update.Message.MessageID, bot)
		case "teams":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowTeams(false, dbConfig), update.Message.MessageID, bot)
		default:
			if strings.HasPrefix(update.Message.Text, "/") {
				_ = src.SendMessageTelegram(update.Message.Chat.ID, "&#9940;I don't know what you want. But you can use /help", update.Message.MessageID, bot)
				continue
			}
			src.CheckCode(update.Message, bot, dbConfig)
		}
	}
}
