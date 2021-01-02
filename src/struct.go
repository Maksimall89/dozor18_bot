package src

import "os"

var NameFileWords = "./words.txt"

type Config struct {
	TelegramBotToken string
	OwnName          string
	ListenPath       string
	Port             string
}

func (conf *Config) Init() {
	if value, exists := os.LookupEnv("TelegramBotToken"); exists {
		conf.TelegramBotToken = value
	}
	if value, exists := os.LookupEnv("OwnName"); exists {
		conf.OwnName = value
	}
	if value, exists := os.LookupEnv("ListenPath"); exists {
		conf.ListenPath = value
	}
	if value, exists := os.LookupEnv("PORT"); exists {
		conf.Port = value
	}
}

type DBconfig struct {
	DBURL        string
	DriverNameDB string
}

func (confDataBase *DBconfig) Init() {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		confDataBase.DBURL = value
	}
	if value, exists := os.LookupEnv("DriverNameDB"); exists {
		confDataBase.DriverNameDB = value
	}
}

type Codes struct {
	ID       int
	Time     string
	NickName string
	Code     string
	Danger   string
	Sector   string
}

type Teams struct {
	ID       int
	Time     string
	NickName string
	Team     string
	Hash     string
}

type Users struct {
	ID       int
	Time     string
	NickName string
	Team     string
	Login    string
	Password string
}
