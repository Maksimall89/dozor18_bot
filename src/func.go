package src

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetNickName(from *tgbotapi.User) string {
	if from.UserName != "" {
		return fmt.Sprintf("%d + %s", from.ID, from.UserName)
	}

	return fmt.Sprintf("%d + %s %s", from.ID, from.FirstName, from.LastName)
}

func GetListHelps(from *tgbotapi.User, adminNickname string) (commandList string) {
	type commandStruct struct {
		admin   bool
		command string
	}

	var commands = []commandStruct{
		{false, "/help - информация по всем доступным командам;\n"},
		{false, "/codes - коды;\n"},
		{false, "/generate, /gen - сгенерировать коды;\n"},
		{false, "/text - текст приквела;\n"},
		{true, "/show - показать все коды;\n"},
		{true, "/reset - удалить все из БД и создать новые;\n"},
		{true, "/add - добавить новые правильные коды в формате: Code,Danger,Sector;\n"},
		{true, "/update - обновить коды в бд, в формате: CodeNew,Danger,Sector,CodeOld;\n"},
		{true, "/delete - удалить указанный код;\n"},
		{true, "/create - создать команду;\n"},
		{true, "/join - вступить в команду;\n"},
		{true, "/list - список участников;\n"},
		{true, "/leave - выйти из команды;\n"},
		{true, "/resetteams - удалить все команды;\n"},
	}

	for _, command := range commands {
		if command.admin && adminNickname != from.UserName {
			continue
		}
		commandList += command.command
	}
	return commandList
}

func SendMessageTelegram(chatId int64, message string, replyToMessageID int, bot *tgbotapi.BotAPI) error {
	var pointerStr int
	var msg tgbotapi.MessageConfig
	var newMsg tgbotapi.Message
	var err error
	isEnd := false

	if len(message) == 0 {
		message = "&#128190;Нет данных."
	}

	if replyToMessageID != 0 {
		msg.ReplyToMessageID = replyToMessageID
	}
	msg.ChatID = chatId
	msg.ParseMode = "HTML"

	for !isEnd {
		if len(message) > 4090 { // ограничение на длину одного сообщения 4096
			pointerStr = strings.LastIndex(message[0:4090], "\n")
			msg.Text = message[0:pointerStr]
			message = message[pointerStr:]
		} else {
			msg.Text = message
			isEnd = true
		}

		newMsg, err = bot.Send(msg)
		if err != nil {
			msg.ParseMode = "Markdown"
			newMsg, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
				log.Println(msg.Text)
			}
			msg.ParseMode = "HTML"
		}
		if strings.Contains(msg.Text, "&#9889;Выдан новый уровень!") {
			_, err := bot.PinChatMessage(tgbotapi.PinChatMessageConfig{ChatID: chatId, MessageID: newMsg.MessageID})
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
