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
	dbConfig := src.DBconfig{}
	dbConfig.Init()

	var str string
	// read from channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.From.UserName == configuration.OwnName {
			switch strings.ToLower(update.Message.Command()) {
			case "reset", "restart":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBResetAllCodes(), 0, bot)
				continue
			case "show":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesAll(dbConfig), 0, bot)
				continue
			case "add":
				// /add code1,1,1
				arrCodes := strings.Split(update.Message.Text, "\n")
				for _, code := range arrCodes {
					code = strings.Replace(code, "/add", "", -1)
					str = dbConfig.DBInsertCodesRight(code)
					_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				}
				continue
			case "update":
				// /update CodeNew,Danger,Sector,CodeOld
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBUpdateCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot)
				continue
			case "delete":
				// /delete code1
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBDeleteCodesRight(update.Message.CommandArguments()), update.Message.MessageID, bot)
				continue
			case "resetteams":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, dbConfig.DBResetUser(), update.Message.MessageID, bot)
				continue
			case "create":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.CreateTeam(update.Message, dbConfig), update.Message.MessageID, bot)
				continue
			case "listall":
				str = "Список всех команд:\n"
				teams := dbConfig.DBSelectTeam("")
				for _, value := range teams {
					str += fmt.Sprintf("%d. <b>%s</b>, %s, <code>%s</code>, %s\n", value.ID, value.Team, value.NickName, value.Hash, value.Time)
				}
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				continue
			case "list":
				// TODO show only my team
				str = "Список всех участников команд:\n"
				users := dbConfig.DBSelectUsers("")
				for key, value := range users {
					str += fmt.Sprintf("%d. <b>%s</b> %s\n", key, value.NickName, value.Team)
				}
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				continue
			case "join":
				_ = src.SendMessageTelegram(update.Message.Chat.ID, src.AddUser(update.Message, dbConfig), update.Message.MessageID, bot)
				continue
			case "leave":
				continue
			}
		}

		switch strings.ToLower(update.Message.Command()) {
		case "codes":
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.ShowCodesMy(update.Message, dbConfig), update.Message.MessageID, bot)
		case "generate", "gen":
			strArr := strings.Split(update.Message.CommandArguments(), ",")
			number, err := strconv.Atoi(strArr[0])
			if err != nil {
				str = "Не по формату:\n<code>/generate 10</code>\n<code>/generate 10,1D,R</code>"
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
			_ = src.SendMessageTelegram(update.Message.Chat.ID, src.GetListHelps(update.Message.From, configuration.OwnName), update.Message.MessageID, bot)
		default:
			if strings.HasPrefix(update.Message.Text, "/") {
				_ = src.SendMessageTelegram(update.Message.Chat.ID, "I don't know what you want. But you can use /help", update.Message.MessageID, bot)
				break
			}
			src.CheckCode(update.Message, bot, dbConfig)
		}
	}
}
