package src

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
	"time"
)

const (
	WhereUserID = "WHERE UserID='%d'"
	WhereTeam   = "WHERE Team='%s'"
)

func CheckCode(message *tgbotapi.Message, bot *tgbotapi.BotAPI, dbConfig Config) {
	dataRight := dbConfig.DBSelectCodesRight()
	codes := Codes{}

	myTeam := ""
	UserID, _ := GetNickName(message.From)
	users := dbConfig.DBSelectUsers(fmt.Sprintf(WhereUserID, UserID))
	if len(users) > 0 {
		myTeam = users[0].Team
	}

	str := "&#9940; Код неверный."
	for _, valueData := range dataRight {
		strArr := strings.Split(valueData.Code, "|")
		for _, value := range strArr {
			if strings.EqualFold(value, strings.TrimSpace(message.Text)) {
				str = fmt.Sprintf("&#9989; Снят код <b>№%d</b> с КО %s из сектора %s", valueData.ID, valueData.Danger, valueData.Sector)
				codes.Time = GetTime()
				codes.UserID, codes.NickName = GetNickName(message.From)
				codes.Code = strings.ToLower(strings.TrimSpace(message.Text))
				codes.Danger = valueData.Danger
				codes.Sector = valueData.Sector
				codes.Team = myTeam
				dbConfig.DBInsertCodesUsers(&codes)
				break
			}
		}
	}
	_ = SendMessageTelegram(message.Chat.ID, str, message.MessageID, bot)
}
func GetInvite(message *tgbotapi.Message, dbConfig Config) string {
	UserID, _ := GetNickName(message.From)
	condition := fmt.Sprintf(WhereUserID, UserID)
	users := dbConfig.DBSelectUsers(condition)
	if len(users) < 1 {
		return "&#10071;Вы не состоите ни в одной команде."
	}
	myTeam := users[0].Team
	condition = fmt.Sprintf(WhereTeam, myTeam)
	teams := dbConfig.DBSelectTeam(condition)
	if len(teams) < 1 {
		return "&#10071;Вы состоите в удаленной команде."
	}
	return fmt.Sprintf("Для вступления в команду <b>%s</b> введите: <code>/join %s, %s </code>", teams[0].Team, teams[0].Team, teams[0].Hash)
}
func ShowCodesAll(dbConfig Config) string {
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
		str += fmt.Sprintf("%d. %s; <b>Ник:</b> %s; <b>Команда:</b> %s; <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s;\n", value.ID, value.Time, value.NickName, value.Team, value.Code, value.Danger, value.Sector)
	}

	return str
}
func ShowCodesMy(message *tgbotapi.Message, dbConfig Config) string {
	var isFound bool
	UserID, _ := GetNickName(message.From)
	condition := fmt.Sprintf(WhereUserID, UserID)
	str := fmt.Sprintf("Коды <b>%s</b>:\n", message.From)
	users := dbConfig.DBSelectUsers(condition)
	if len(users) > 0 {
		condition = fmt.Sprintf(WhereTeam, users[0].Team)
		str = fmt.Sprintf("Коды команды <b>%s</b>:\n", users[0].Team)
	}
	dataAll := dbConfig.DBSelectCodesUser(condition)
	dataRight := dbConfig.DBSelectCodesRight()

	for _, valueData := range dataRight {
		strArr := strings.Split(valueData.Code, "|")
		isFound = false
		for _, value := range strArr {
			for _, base := range dataAll {
				if strings.ToLower(strings.TrimSpace(value)) == base.Code {
					str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор <b>%s</b>, &#9989;<b>СНЯТ</b> (%s)\n", valueData.ID, valueData.Danger, valueData.Sector, valueData.Code)
					isFound = true
					break
				}
			}
		}
		if !isFound {
			str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор: <b>%s</b>, &#10060;<b>НЕ</b> снят\n", valueData.ID, valueData.Danger, valueData.Sector)
		}
	}
	return str
}
func CreateTeam(message *tgbotapi.Message, dbConfig Config) string {
	if len(message.CommandArguments()) < 3 {
		return "&#10071;Слишком короткое название команды, надо минимум 3 символа."
	}
	err := CheckMessage(message.CommandArguments())
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	team := Teams{}
	team.Time = GetTime()
	team.UserID, team.NickName = GetNickName(message.From)
	team.Team = strings.ToLower(strings.TrimSpace(message.CommandArguments()))
	team.Hash = GetMD5Hash(team.Team)

	return dbConfig.DBInsertTeam(&team)
}
func JoinTeam(message *tgbotapi.Message, dbConfig Config) string {
	strArr := strings.Split(message.CommandArguments(), ",")
	if len(strArr) < 2 {
		return "&#10071;Нет всех аргументов: /join team, secret key"
	}
	err := CheckMessage(strArr[0])
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	for number, value := range strArr {
		strArr[number] = strings.ToLower(strings.TrimSpace(value))
	}
	user := Users{}
	team := dbConfig.DBSelectTeam(fmt.Sprintf(WhereTeam, strArr[0]))
	if len(team) != 1 || strArr[1] != team[0].Hash {
		return "&#10071;Неверный ключ или название команды"
	}
	user.UserID, user.NickName = GetNickName(message.From)
	user.Time = GetTime()
	user.Team = strArr[0]

	return dbConfig.DBInsertUser(&user)
}
func ShowUsers(message *tgbotapi.Message, isMyTeam bool, dbConfig Config) string {
	var users []Users
	var condition string
	str := "Список всех участников в командах:\n"
	UserID, _ := GetNickName(message.From)
	if isMyTeam {
		condition = fmt.Sprintf(WhereUserID, UserID)
		users = dbConfig.DBSelectUsers(condition)
		if len(users) < 1 {
			return "&#10071;Вы не состоите ни в одной команде."
		}
		condition = fmt.Sprintf(WhereTeam, users[0].Team)
		str = fmt.Sprintf("Список всех участников команды <b>%s</b>:\n", users[0].Team)
	}
	users = dbConfig.DBSelectUsers(condition)
	for key, value := range users {
		str += fmt.Sprintf("%d. <b>%s</b>; %s\n", key, value.NickName, value.Team)
	}
	return str
}
func GetTime() string {
	return fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond())
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	text = hex.EncodeToString(hasher.Sum(nil))
	return text[:18]
}
func GetNickName(from *tgbotapi.User) (int, string) {
	if from.UserName != "" {
		return from.ID, from.UserName
	}
	return from.ID, fmt.Sprintf("%s %s", from.FirstName, from.LastName)
}
func GetListHelps(from *tgbotapi.User, adminID int) (commandList string) {
	type commandStruct struct {
		admin   bool
		command string
	}

	var commands = []commandStruct{
		{false, "/help - информация по всем доступным командам;\n"},
		{false, "/codes - коды;\n"},
		{false, "/generate, /gen - сгенерировать коды;\n"},
		{false, "/text - текст приквела;\n"},
		{false, "/create - создать команду;\n"},
		{false, "/join - вступить в команду;\n"},
		{false, "/list - список участников команды;\n"},
		{false, "/listusers - список участников в командах;\n"},
		{false, "/leave - выйти из команды;\n"},
		{false, "/invite - получить ссылку приглашение в команду;\n"},
		{true, "/show - показать все коды;\n"},
		{true, "/reset - удалить все из БД и создать новые;\n"},
		{true, "/add - добавить новые правильные коды в формате: Code,Danger,Sector;\n"},
		{true, "/update - обновить коды в бд, в формате: CodeNew,Danger,Sector,CodeOld;\n"},
		{true, "/delete - удалить указанный код;\n"},
		{true, "/listteams - список всех команд;\n"},
		{true, "/resetteams - удалить все команды;\n"},
	}

	for _, command := range commands {
		if command.admin && adminID != from.ID {
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
func CheckMessage(message string) error {
	if strings.ContainsAny(strings.ToLower(message), "\"`~-\\=:;/,.'*+@#№%$%^&(){}[]|") {
		return errors.New("&#10071;Недопустимые символы в сообщении. Можно использовать лишь буквы и цифры русского и английского алфавита.")
	}
	return nil
}
