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
	if value, exists := os.LookupEnv("Port"); exists {
		conf.Port = value
	}
}

type DataBase struct {
	Number       int
	Time         string
	NickName     string
	Code         string
	Danger       string
	Sector       string
	DBURL        string
	DriverNameDB string
}

func (confDataBase *DataBase) Init() {
	if value, exists := os.LookupEnv("DATABASE_URL"); exists {
		confDataBase.DBURL = value
	}
	if value, exists := os.LookupEnv("DriverNameDB"); exists {
		confDataBase.DriverNameDB = value
	}
}

func (newConf *Teams) InitDB(oldConf DataBase) {
	newConf.DriverNameDB = oldConf.DriverNameDB
	newConf.DBURL = oldConf.DBURL
}

type Teams struct {
	ID           int
	Time         string
	Team         string
	Hash         string
	Owner        string
	DBURL        string
	DriverNameDB string
}
type Users struct {
	ID           int
	NickName     string
	Time         string
	Team         string
	Login        string
	Password     string
	DBURL        string
	DriverNameDB string
}
