package src

import (
	"os"
	"strconv"
)

const NameFileWords = "./words.txt"

type Config struct {
	TelegramBotToken string
	OwnName          string
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
	if value, exists := os.LookupEnv("OwnName"); exists {
		dbConfig.OwnName = value
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
