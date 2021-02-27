package src

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
)

const (
	whereUserID = "WHERE UserID='%d'"
	whereTeam   = "WHERE Team='%s'"
)

func ShowTeams(isAllInfo bool, dbConfig Config) string {
	str := "&#9745;Список всех команд:\n"
	teams := dbConfig.DBSelectTeam("")
	for number, value := range teams {
		if value.Team == "void" {
			continue
		}
		str += fmt.Sprintf("\n%d. Команда: <b>%s</b>; Капитан: <b>%s</b>", number+1, value.Team, value.NickName)
		if isAllInfo {
			str += fmt.Sprintf(" <code>%s</code>, %s", value.Hash, value.Time)
		}
	}
	return str
}
func CheckCode(message *tgbotapi.Message, bot *tgbotapi.BotAPI, dbConfig Config) {
	var codes Codes
	var myTeam string

	UserID, NickName := GetNickName(message.From)
	users := dbConfig.DBSelectUsers(fmt.Sprintf(whereUserID, UserID))
	if len(users) > 0 {
		myTeam = users[0].Team
	} else {
		myTeam = "void"
	}

	code := strings.ToLower(strings.TrimSpace(message.Text))
	str := "&#10060; Код неверный."
	dataRight := dbConfig.DBSelectCodesRight()
	for numberRight, valueRight := range dataRight {
		strArr := strings.Split(valueRight.Code, "|")
		for _, value := range strArr {
			if value == code {
				str = fmt.Sprintf("&#9989;Снят код <b>№%d</b> с КО <b>%s</b> из сектора <b>%s</b>", numberRight+1, valueRight.Danger, valueRight.Sector)
				codes.UserID = UserID
				codes.NickName = NickName
				codes.Code = valueRight.Code
				codes.Team = myTeam
				dbConfig.DBInsertCodesUsers(&codes)
				break
			}
		}
	}
	_ = SendMessageTelegram(message.Chat.ID, str, message.MessageID, bot, "main")
}
func GetInvite(message *tgbotapi.Message, dbConfig Config) string {
	str := "&#10071;Вы не состоите ни в одной команде."
	UserID, _ := GetNickName(message.From)
	condition := fmt.Sprintf(whereUserID, UserID)
	users := dbConfig.DBSelectUsers(condition)
	if len(users) < 1 {
		return str
	}
	if users[0].Team == "void" {
		return str
	}
	myTeam := users[0].Team
	condition = fmt.Sprintf(whereTeam, myTeam)
	teams := dbConfig.DBSelectTeam(condition)
	if len(teams) < 1 {
		return "&#10071;Вы состоите в удаленной команде."
	}
	return fmt.Sprintf("&#9745;Для вступления в команду <b>%s</b> введите: <code>/join %s, %s </code>", teams[0].Team, teams[0].Team, teams[0].Hash)
}
func ShowCodesAll(dbConfig Config) string {
	dataAllRight := dbConfig.DBSelectCodesRight()
	// ID, Time, NickName, Code, Danger, Sector
	str := fmt.Sprintf("Всего кодов в движке: <b>%d</b>\n&#9989;Коды верные:\n", len(dataAllRight))
	for number, value := range dataAllRight {
		str += fmt.Sprintf("%d. <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s; <b>Бонус:</b> %d сек;\n", number+1, value.Code, value.Danger, value.Sector, value.TimeBonus)
	}

	dataAllUsers := dbConfig.DBSelectCodesUser("")
	// ID, Time, NickName, Code, Danger, Sector
	str += fmt.Sprintf("\nВсего кодов введено: <b>%d</b>\n&#9745;Коды Юзеров:\n", len(dataAllUsers))
	for number, value := range dataAllUsers {
		str += fmt.Sprintf("%d. %s; <b>Ник:</b> %s; <b>Команда:</b> %s; <b>Код:</b> %s;\n", number+1, value.Time, value.NickName, value.Team, value.Code)
	}

	return str
}
func ShowCodesMy(message *tgbotapi.Message, dbConfig Config) string {
	var isFound bool
	UserID, _ := GetNickName(message.From)
	condition := fmt.Sprintf(whereUserID, UserID)
	str := fmt.Sprintf("&#9745;Коды <b>%s</b>:\n", message.From)
	users := dbConfig.DBSelectUsers(condition)
	if len(users) > 0 {
		if users[0].Team != "void" {
			condition = fmt.Sprintf(whereTeam, users[0].Team)
			str = fmt.Sprintf("&#9745;Коды команды <b>%s</b>:\n", users[0].Team)
		}
	}
	dataUser := dbConfig.DBSelectCodesUser(condition)
	dataRight := dbConfig.DBSelectCodesRight()

	for number, valueRight := range dataRight {
		isFound = false
		for _, valueUser := range dataUser {
			if strings.ToLower(strings.TrimSpace(valueUser.Code)) == valueRight.Code {
				str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор <b>%s</b>, &#9989;<b>СНЯТ</b> (%s), бонус <b>%d</b> сек\n", number+1, valueRight.Danger, valueRight.Sector, valueRight.Code, valueRight.TimeBonus)
				isFound = true
				break
			}
		}
		if !isFound {
			str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор: <b>%s</b>, &#10060;<b>НЕ</b> снят, бонус <b>%d</b> сек\n", number+1, valueRight.Danger, valueRight.Sector, valueRight.TimeBonus)
		}
	}
	return str
}
func CreateTeam(message *tgbotapi.Message, dbConfig Config) string {
	err := CheckMessage(message.CommandArguments())
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	team := Teams{}
	team.UserID, team.NickName = GetNickName(message.From)
	team.Team = strings.ToLower(strings.TrimSpace(message.CommandArguments()))
	team.Hash = GetMD5Hash(team.Team, dbConfig.ListenPath)

	return dbConfig.DBInsertTeam(&team)
}
func JoinTeam(message *tgbotapi.Message, dbConfig Config) string {
	strArr := strings.Split(message.CommandArguments(), ",")
	if len(strArr) != 2 {
		return "&#10071;Нет всех аргументов: <code>/join team, secret key</code>"
	}
	for number, value := range strArr {
		err := CheckMessage(value)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}
		strArr[number] = strings.ToLower(strings.TrimSpace(value))
	}
	user := Users{}
	team := dbConfig.DBSelectTeam(fmt.Sprintf(whereTeam, strArr[0]))
	if len(team) != 1 || strArr[1] != team[0].Hash {
		return "&#10071;Неверный ключ или название команды"
	}
	user.UserID, user.NickName = GetNickName(message.From)
	user.Team = strArr[0]

	return dbConfig.DBInsertUser(&user)
}
func ShowUsers(message *tgbotapi.Message, isMyTeam bool, dbConfig Config) string {
	var users []Users
	var condition string
	str := "&#9745;Список всех участников в командах:\n"
	UserID, _ := GetNickName(message.From)
	if isMyTeam {
		condition = fmt.Sprintf(whereUserID, UserID)
		users = dbConfig.DBSelectUsers(condition)
		if len(users) < 1 {
			return "&#10071;Вы не состоите ни в одной команде."
		}
		condition = fmt.Sprintf(whereTeam, users[0].Team)
		str = fmt.Sprintf("&#9745;Список всех участников команды <b>%s</b>:\n", users[0].Team)
	}
	users = dbConfig.DBSelectUsers(condition)
	for key, value := range users {
		str += fmt.Sprintf("%d. Ник: <b>%s</b>; Команда: <b>%s</b>\n", key+1, value.NickName, value.Team)
	}
	return str
}

func GetMD5Hash(text string, key string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(text + key))
	text = hex.EncodeToString(hash.Sum(nil))
	return text[:18]
}
func GetNickName(from *tgbotapi.User) (int, string) {
	if from.UserName != "" {
		return from.ID, from.UserName
	}
	return from.ID, fmt.Sprintf("%s %s", from.FirstName, from.LastName)
}
func GetListHelps(from *tgbotapi.User, adminID int) (commandList string) {
	for _, command := range Commands {
		if command.IsAdmin && adminID != from.ID {
			continue
		}
		commandList += fmt.Sprintf("/%s - %s;\n", command.Command, command.Describe)
	}
	return commandList
}
func SendMessageTelegram(chatId int64, message string, replyToMessageID int, bot *tgbotapi.BotAPI, levelButtons string) error {
	var pointerStr int
	var msg tgbotapi.MessageConfig
	var err error
	isEnd := false

	if len(message) == 0 {
		message = "&#9940;Нет данных."
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, button := range Commands {
		if button.LevelMenu != levelButtons && button.LevelMenu != "all" {
			continue
		}
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(button.Describe, button.Command)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	msg.ReplyMarkup = keyboard

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

		_, err = bot.Send(msg)
		if err != nil {
			msg.ParseMode = "Markdown"
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
				log.Println(msg.Text)
			}
			msg.ParseMode = "HTML"
		}
	}
	return nil
}
func CheckMessage(message string) error {
	if len(message) > 100 {
		return errors.New("&#10071;Сообщение слишком длинное")
	}
	if len(message) < 3 {
		return errors.New("&#10071;Сообщение слишком короткое")
	}
	if ok, _ := regexp.MatchString(`^[а-яА-ЯёЁa-zA-Z0-9 ]+$`, message); !ok {
		return errors.New("&#10071;Недопустимые символы в сообщении. Можно использовать лишь буквы и цифры русского и английского алфавита")
	}
	return nil
}
