package src

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strconv"
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
	_ = SendMessageTelegram(message.Chat.ID, str, message.MessageID, bot, "codes")
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
		str += fmt.Sprintf("%d. <b>Код:</b> %s; <b>КО:</b> %s; <b>Сектор:</b> %s; <b>Бонус:</b> %s;\n", number+1, value.Code, value.Danger, value.Sector, ConvertTimeSec(value.TimeBonus))
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
			if valueUser.Code == valueRight.Code {
				str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор <b>%s</b>, &#9989;<b>СНЯТ</b> (%s), бонус <b>%s</b>, задание <b>%d</b>\n", number+1, valueRight.Danger, valueRight.Sector, valueRight.Code, ConvertTimeSec(valueRight.TimeBonus), valueRight.TaskID)
				isFound = true
				break
			}
		}
		if !isFound {
			str += fmt.Sprintf("%d. КО: <b>%s</b>, сектор: <b>%s</b>, &#10060;<b>НЕ</b> снят, бонус <b>%s</b>, задание <b>%d</b>\n", number+1, valueRight.Danger, valueRight.Sector, ConvertTimeSec(valueRight.TimeBonus), valueRight.TaskID)
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
	var pointer int
	var msg tgbotapi.MessageConfig
	var err error
	isEnd := false

	if len(message) == 0 {
		message = "&#9940;Нет данных."
	}

	if replyToMessageID != 0 {
		msg.ReplyToMessageID = replyToMessageID
	}
	msg.ChatID = chatId
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = createKeyboard(levelButtons)
	for !isEnd {
		if len(message) > 4090 { // ограничение на длину одного сообщения 4096
			pointer = strings.LastIndex(message[0:4090], "\n")
			msg.Text = message[0:pointer]
			message = message[pointer:]
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
func ArrTrimSpace(arr []string) {
	for number, value := range arr {
		arr[number] = strings.TrimSpace(value)
	}
}
func ConvertTimeSec(times int) string {
	if times == 0 {
		return "0 секунд"
	}
	// из секунд в минуту
	str := ""
	timeSec := times % 60
	timeMin := times / 60
	timeHour := times / 3600
	if timeHour > 0 {
		timeMin = timeMin % 60
	}
	timeDay := times / 86400
	if timeDay > 0 {
		timeMin = times % 60
		timeHour = timeHour % 24
	}

	// Дни
	switch timeDay {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d день ", timeDay)
	case 2, 3, 4:
		str += fmt.Sprintf("%d дня ", timeDay)
	default:
		str += fmt.Sprintf("%d дней ", timeDay)
	}
	// Часы
	switch timeHour {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d час ", timeHour)
	case 2, 3, 4:
		str += fmt.Sprintf("%d часа ", timeHour)
	default:
		str += fmt.Sprintf("%d часов ", timeHour)
	}
	// Минуты
	switch timeMin {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d минута ", timeMin)
	case 2, 3, 4:
		str += fmt.Sprintf("%d минуты ", timeMin)
	default:
		str += fmt.Sprintf("%d минут ", timeMin)
	}
	// Секунды
	switch timeSec {
	case 0:
		str += ""
	case 1:
		str += fmt.Sprintf("%d секунда", timeSec)
	case 2, 3, 4:
		str += fmt.Sprintf("%d секунды", timeSec)
	default:
		str += fmt.Sprintf("%d секунд", timeSec)
	}
	return strings.TrimSpace(str)
}
func createKeyboard(levelButtons string) (keyboard tgbotapi.InlineKeyboardMarkup) {
	var counter int
	var row []tgbotapi.InlineKeyboardButton
	for _, command := range Commands {
		for _, levelMenu := range command.LevelMenu {
			if levelMenu != levelButtons && levelMenu != "all" {
				continue
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(command.Describe, command.Command))
			counter++
			if counter%3 == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				row = nil
			}
		}
	}
	if counter < 3 || counter == len(Commands) {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return keyboard
}
func GetTasks(message *tgbotapi.Message, dbConfig Config) string {
	idMap := make(map[int]int)
	var counter int
	var err error
	var str string

	numberArr := strings.Split(message.CommandArguments(), " ")
	for _, number := range numberArr {
		idMap[counter], err = strconv.Atoi(number)
		if err == nil {
			counter++
		}
	}

	if counter == 0 {
		str = createTaskList(dbConfig, "", 1)
	} else {
		for _, idTask := range idMap {
			str += createTaskList(dbConfig, fmt.Sprintf("WHERE ID='%d'", idTask), 1)
		}
	}

	if len(str) == 0 {
		return `Текст приквела доступен на нашем сайте <a href="http://dozor18.ru">http://dozor18.ru</a>.`

	}
	return str
}
func createTaskList(dbConfig Config, condition string, repeat int) (taskString string) {
	tasks := dbConfig.DBSelectTask(condition)
	for _, task := range tasks {
		taskString += fmt.Sprintf("<b>%d</b>. %s\n", task.ID, task.Text)
	}
	if taskString == "" && repeat == 1 {
		repeat = 0
		taskString = createTaskList(dbConfig, "", repeat)
	}
	return taskString
}
func CreateTask(message *tgbotapi.Message, dbConfig Config) string {
	strArr := strings.Split(message.CommandArguments(), ",")
	if len(strArr) != 2 {
		return "&#10071;Нет всех аргументов: <code>/addtask id, text task</code>"
	}
	id, err := strconv.Atoi(strArr[0])
	if err != nil {
		return "&#10071;Id передан неверно"
	}
	task := Tasks{}
	task.ID = id
	task.Text = strArr[1]
	return dbConfig.DBInsertTask(&task)
}
