package src

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
	"time"
)

func CheckCode(user *tgbotapi.Message, bot *tgbotapi.BotAPI, dbConfig DBconfig) {
	dataRight := dbConfig.DBSelectCodesRight()
	codes := Codes{}
	str := ""

	for _, valueData := range dataRight {
		strArr := strings.Split(valueData.Code, "|")
		str = "&#9940; Код неверный."
		for _, value := range strArr {
			if strings.EqualFold(value, strings.TrimSpace(user.Text)) {
				str = fmt.Sprintf("&#9989; Снят код <b>№%d</b> с КО %s из сектора %s", valueData.ID, valueData.Danger, valueData.Sector)
				codes.Time = GetTime()
				codes.NickName = GetNickName(user.From)
				codes.Code = strings.ToLower(strings.TrimSpace(user.Text))
				codes.Danger = valueData.Danger
				codes.Sector = valueData.Sector
				dbConfig.DBInsertCodesUsers(&codes)
				break
			}
		}
		_ = SendMessageTelegram(user.Chat.ID, str, user.MessageID, bot)
	}
}
func ShowCodesAll(dbConfig DBconfig) string {
	dataAllRight := dbConfig.DBSelectCodesRight()
	// ID, Time, NickName, Code, Danger, Sector
	str := fmt.Sprintf("Всего кодов в движке: %d\n&#9989;Коды верные:\n", len(dataAllRight))
	for _, value := range dataAllRight {
		str += fmt.Sprintf("%d. <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s;\n", value.ID, value.Code, value.Danger, value.Sector)
	}

	dataAllUsers := dbConfig.DBSelectCodesUser("")
	// ID, Time, NickName, Code, Danger, Sector
	str += fmt.Sprintf("\nВсего кодов введено: %d\n&#9745;Коды Юзеров:\n", len(dataAllUsers))
	for _, value := range dataAllUsers {
		str += fmt.Sprintf("%d. %s; <b>Ник:</b> %s; <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s;\n", value.ID, value.Time, value.NickName, value.Code, value.Danger, value.Sector)
	}

	return str
}
func ShowCodesMy(user *tgbotapi.Message, dbConfig DBconfig) string {
	var isFound bool
	condition := fmt.Sprintf("WHERE NickName = '%s'", GetNickName(user.From))
	str := fmt.Sprintf("Коды <b>%s</b>:\n", user.From)
	team := dbConfig.DBSelectTeam(GetNickName(user.From))
	if len(team) > 0 {
		condition = fmt.Sprintf("WHERE Team = '%s'", team[0].Team)
		str = fmt.Sprintf("Коды команды <b>%s</b>:\n", team[0].Team)
	}
	dataAll := dbConfig.DBSelectCodesUser(condition)
	dataRight := dbConfig.DBSelectCodesRight()

	for _, valueData := range dataRight {
		strArr := strings.Split(valueData.Code, "|")
		for _, value := range strArr {
			isFound = false
			for _, base := range dataAll {
				if strings.ToLower(strings.TrimSpace(value)) == base.Code {
					str += fmt.Sprintf("%d. Код Опасности: <b>%s</b>, сектор <b>%s</b>, &#9989;<b>СНЯТ</b>\n", valueData.ID, valueData.Danger, valueData.Sector)
					isFound = true
					break
				}
			}
			if !isFound {
				str += fmt.Sprintf("%d. Код Опасности: <b>%s</b>, сектор: <b>%s</b>, &#10060;<b>НЕ</b> снят\n", valueData.ID, valueData.Danger, valueData.Sector)
			}
		}
	}
	return str
}
func CreateTeam(message *tgbotapi.Message, dbConfig DBconfig) string {
	if len(message.CommandArguments()) < 3 {
		return "&#10071;Слишком короткое название команды, надо минимум 3 символа."
	}

	team := Teams{}
	team.Time = GetTime()
	team.NickName = GetNickName(message.From)
	team.Team = strings.ToLower(strings.TrimSpace(message.CommandArguments()))
	team.Hash = GetMD5Hash(team.Team)

	teams := dbConfig.DBSelectTeam("")
	for _, value := range teams {
		if value.Team == team.Team {
			return "&#10071; Такая команда уже есть."
		}
	}
	return dbConfig.DBCreateTeam(&team)
}
func JoinTeam(addUser *tgbotapi.Message, dbConfig DBconfig) string {
	strArr := strings.Split(addUser.CommandArguments(), ",")
	if len(strArr) < 2 {
		return "&#10071;Нет всех аргументов: /join team, secret key"
	}
	for number, value := range strArr {
		strArr[number] = strings.ToLower(strings.TrimSpace(value))
	}
	user := Users{}
	team := dbConfig.DBSelectTeam(strArr[0])
	if len(team) != 1 || strArr[1] != team[0].Hash {
		return "&#10071;Неверный ключ или имя команды"
	}
	user.NickName = GetNickName(addUser.From)
	user.Time = GetTime()
	user.Team = strArr[0]

	return dbConfig.DBInsertUser(&user)
}
func ShowUsers(team *tgbotapi.Message, isAdmin bool, dbConfig DBconfig) string {
	if len(team.CommandArguments()) < 3 && !isAdmin {
		return "&#10071;Слишком короткое имя команды!"
	}
	condition := fmt.Sprintf("WHERE Team = '%s'", strings.ToLower(team.CommandArguments()))
	str := fmt.Sprintf("Список всех участников команды <b>%s</b>:\n", team.CommandArguments())
	users := dbConfig.DBSelectUsers(condition)
	for key, value := range users {
		str += fmt.Sprintf("%d. <b>%s</b> %s\n", key, value.NickName, value.Team)
	}
	return str
}
func GetTime() string {
	return fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	text = hex.EncodeToString(hasher.Sum(nil))
	return text[:18]
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
		{true, "/list - список участников команды;\n"},
		{true, "/listusers - список участников команд;\n"},
		{true, "/listteams - список всех команд;\n"},
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
