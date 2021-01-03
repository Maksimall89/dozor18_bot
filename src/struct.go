package src

import (
	"os"
	"strconv"
)

var NameFileWords = "./words.txt"

type Config struct {
	TelegramBotToken string
	OwnName          string
	OwnID            int
	ListenPath       string
	Port             string
	DBURL            string
	DriverNameDB     string
}

func (conf *Config) Init() {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		conf.DBURL = value
	}
	if value, exists := os.LookupEnv("DriverNameDB"); exists {
		conf.DriverNameDB = value
	}
	if value, exists := os.LookupEnv("TelegramBotToken"); exists {
		conf.TelegramBotToken = value
	}
	if value, exists := os.LookupEnv("OwnName"); exists {
		conf.OwnName = value
	}
	if value, exists := os.LookupEnv("UserID"); exists {
		conf.OwnID, _ = strconv.Atoi(value)
	}
	if value, exists := os.LookupEnv("ListenPath"); exists {
		conf.ListenPath = value
	}
	if value, exists := os.LookupEnv("PORT"); exists {
		conf.Port = value
	}
}

type Codes struct {
	ID       int
	Time     string
	NickName string
	UserID   int
	Code     string
	Danger   string
	Team     string
	Sector   string
}

type Teams struct {
	ID       int
	Time     string
	NickName string
	UserID   int
	Team     string
	Hash     string
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
