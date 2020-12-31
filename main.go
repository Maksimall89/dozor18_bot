package main

import (
	"dozor18_bot/src"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	data := src.DataBase{}
	data.Init()

	var str string
	// read from channel
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.From.UserName == configuration.OwnName {
			switch strings.ToLower(update.Message.Command()) {
			case "reset", "restart":
				str = data.DBResetAll()
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, 0, bot)
				continue
			case "show":
				dataAll := data.DBSelectAllCodesRight()
				// Number, Time, NickName, Code, Danger, Sector
				str = fmt.Sprintf("Всего кодов в движке: %d\n&#9989;Коды верные:\n", len(dataAll))
				for _, value := range dataAll {
					str += fmt.Sprintf("%d. <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s;\n", value.Number, value.Code, value.Danger, value.Sector)
				}

				dataAll = data.DBSelectAllCodesUser()
				// Number, Time, NickName, Code, Danger, Sector
				str += fmt.Sprintf("\nВсего кодов введено: %d\n&#9745;Коды Юзеров:\n", len(dataAll))
				for _, value := range dataAll {
					str += fmt.Sprintf("%d. %s; <b>Ник:</b> %s; <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s;\n", value.Number, value.Time, value.NickName, value.Code, value.Danger, value.Sector)
				}
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, 0, bot)
				continue
			case "add":
				// /add code1,1,1
				arrCodes := strings.Split(update.Message.Text, "\n")
				for _, code := range arrCodes {
					code = strings.Replace(code, "/add", "", -1)
					str = data.DBInsertCodesRight(code)
					_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				}
				continue
			case "update":
				// /update CodeNew,Danger,Sector,CodeOld
				str = data.DBUpdateCodesRight(update.Message.CommandArguments())
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				continue
			case "delete":
				// /delete code1
				str = data.DBDeleteCodesRight(update.Message.CommandArguments())
				_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
				continue
			}
		}

		switch strings.ToLower(update.Message.Command()) {
		case "codes":
			var isFound bool
			str = ""

			data.NickName = src.GetNickName(update.Message.From)
			dataAll := data.DBSelectCodes()
			dataRight := data.DBSelectAllCodesRight()

			for _, valueData := range dataRight {
				strArr := strings.Split(valueData.Code, "|")
				for _, value := range strArr {
					isFound = false
					for _, base := range dataAll {
						if strings.ToLower(strings.TrimSpace(value)) == base.Code {
							str += fmt.Sprintf("%d. Код Опасности: <b>%s</b>, сектор <b>%s</b>, &#9989;<b>СНЯТ</b>\n", valueData.Number, valueData.Danger, valueData.Sector)
							isFound = true
							break
						}
					}

					if !isFound {
						str += fmt.Sprintf("%d. Код Опасности: <b>%s</b>, сектор: <b>%s</b>, &#10060;<b>НЕ</b> снят\n", valueData.Number, valueData.Danger, valueData.Sector)
					}
				}
			}
			_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)
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

			str = ""
			dataRight := data.DBSelectAllCodesRight()
			for _, valueData := range dataRight {
				strArr := strings.Split(valueData.Code, "|")
				for _, value := range strArr {
					if strings.EqualFold(value, strings.TrimSpace(update.Message.Text)) {
						str = fmt.Sprintf("&#9989; Снят код <b>№%d</b> с КО %s из сектора %s", valueData.Number, valueData.Danger, valueData.Sector)
						_ = src.SendMessageTelegram(update.Message.Chat.ID, str, update.Message.MessageID, bot)

						data.Time = fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
						data.NickName = src.GetNickName(update.Message.From)
						data.Code = strings.ToLower(strings.TrimSpace(update.Message.Text))
						data.Danger = valueData.Danger
						data.Sector = valueData.Sector
						data.DBInsertCodesUsers()
						break
					}
				}
			}

			if str == "" {
				_ = src.SendMessageTelegram(update.Message.Chat.ID, "&#9940; Код неверный.", update.Message.MessageID, bot)
			}
		}
	}
}
