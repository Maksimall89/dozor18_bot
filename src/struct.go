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
	LevelMenu []string
}

var Commands = []CommandStruct{
	{false, "все команды", "help", []string{"all"}},
	{false, "список кодов", "codes", []string{"codes", "main"}},
	{false, "сгенерировать коды", "gen", []string{""}},
	{false, "текст задания", "text", []string{"main"}},
	{false, "создать команду", "create имя команды", []string{""}},
	{false, "вступить в команду", "join", []string{""}},
	{false, "список участников команды", "list", []string{"team"}},
	{false, "список участников в командах", "listusers", []string{"team"}},
	{false, "выйти из команды", "leave", []string{"team"}},
	{false, "получить ссылку приглашение в команду", "invite", []string{"team"}},
	{false, "список всех команд", "teams", []string{"team"}},
	{true, "показать все коды", "show", []string{"admin"}},
	{true, "удалить данные из таблицы <b>teams</b> или <b>codes</b>", "reset", []string{""}},
	{true, "добавить новые правильные коды в формате: <b>Code,Danger,Sector,TimeBonus,Tasks</b>", "add", []string{""}},
	{true, "обновить коды в бд, в формате: <b>CodeNew,Danger,Sector,TimeBonus,TaskID,CodeOld</b>", "update", []string{""}},
	{true, "удалить указанный код", "delete", []string{""}},
	{true, "список всех команд", "listteams", []string{"admin"}},
	{true, "создать таблицы в БД", "createdb", []string{""}},
}
