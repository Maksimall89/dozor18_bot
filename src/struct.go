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
	{false, "участники команды", "list", []string{"team"}},
	{false, "участники в командах", "listusers", []string{"team"}},
	{false, "выйти из команды", "leave", []string{"team"}},
	{false, "приглашение в команду", "invite", []string{"team"}},
	{false, "список команд", "teams", []string{"team"}},
	{true, "показать коды", "show", []string{"admin"}},
	{true, "удалить данные из <b>teams</b> или <b>codes</b>", "reset", []string{"admin"}},
	{true, "добавить код: <b>Code,Danger,Sector,TimeBonus,Tasks</b>", "add", []string{"admin"}},
	{true, "обновить коды: <b>CodeNew,Danger,Sector,TimeBonus,TaskID,CodeOld</b>", "update", []string{"admin"}},
	{true, "удалить код", "delete", []string{"admin"}},
	{true, "список команд", "listteams", []string{"admin"}},
	{true, "создать в БД", "createdb", []string{"admin"}},
	{true, "создать задание", "createtask", []string{"admin"}},
}
