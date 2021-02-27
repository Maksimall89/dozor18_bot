package src

import (
	"os"
	"strconv"
)

const NameFileWords = "./words.txt"

type Config struct {
	TelegramBotToken string
	OwnID            int
	ListenPath       string
	Port             string
	DBURL            string
	DriverNameDB     string
}

func (dbConfig *Config) Init() {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		dbConfig.DBURL = value
	}
	if value, exists := os.LookupEnv("DriverNameDB"); exists {
		dbConfig.DriverNameDB = value
	}
	if value, exists := os.LookupEnv("TelegramBotToken"); exists {
		dbConfig.TelegramBotToken = value
	}
	if value, exists := os.LookupEnv("OwnID"); exists {
		dbConfig.OwnID, _ = strconv.Atoi(value)
	}
	if value, exists := os.LookupEnv("ListenPath"); exists {
		dbConfig.ListenPath = value
	}
	if value, exists := os.LookupEnv("PORT"); exists {
		dbConfig.Port = value
	}
}

type Codes struct {
	ID        int
	Time      string
	NickName  string
	UserID    int
	Code      string
	Danger    string
	Team      string
	Sector    string
	TimeBonus int
	TaskID    int
}

type Teams struct {
	ID       int
	Time     string
	NickName string
	UserID   int
	Team     string
	Hash     string
}

type Tasks struct {
	ID   int
	Text string
}

type Users struct {
	ID       int
	Time     string
	NickName string
	UserID   int
	Team     string
	Login    string
	Password string
}

type CommandStruct struct {
	IsAdmin   bool
	Describe  string
	Command   string
	LevelMenu string
}

var Commands = []CommandStruct{
	{false, "информация по всем доступным командам", "help", "all"},
	{false, "список кодов", "codes", "main"},
	{false, "сгенерировать коды", "gen", ""},
	{false, "текст приквела", "text", "main"},
	{false, "создать команду", "create имя команды", ""},
	{false, "вступить в команду", "join", ""},
	{false, "список участников команды", "list", "team"},
	{false, "список участников в командах", "listusers", "team"},
	{false, "выйти из команды", "leave", "team"},
	{false, "получить ссылку приглашение в команду", "invite", "team"},
	{false, "список всех команд", "teams", "team"},
	{true, "показать все коды", "show", "admin"},
	{true, "удалить данные из таблицы <b>teams</b> или <b>codes</b>", "reset", ""},
	{true, "добавить новые правильные коды в формате: <b>Code,Danger,Sector</b>", "add", ""},
	{true, "обновить коды в бд, в формате: <b>CodeNew,Danger,Sector,CodeOld</b>", "update", ""},
	{true, "удалить указанный код", "delete", ""},
	{true, "список всех команд", "listteams", "admin"},
	{true, "создать таблицы в БД", "createdb", ""},
}
