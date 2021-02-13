package src

//import (
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"testing"
//)
//
//func TestDBCreateTables(t *testing.T) {
//	t.Parallel()
//
//	dbConfig := Config{}
//	dbConfig.Init()
//
//	dbConfig.DBCreateTables()
//
//	for _, pair := range tests {
//		resultID, resultStr := GetNickName(pair.input)
//		if resultStr != pair.outputStr || resultID != pair.outputInt {
//			t.Errorf("For %v\nexpected %s %d\ngot %s %d", pair.input, pair.outputStr, pair.outputInt, resultStr, resultID)
//		}
//	}
//}
//
//func TestDBCreateTables(t *testing.T) {
//	t.Parallel()
//
//	type testPair struct {
//		input     *tgbotapi.User
//		outputStr string
//		outputInt int
//	}
//
//	var tests = []testPair{
//		{&tgbotapi.User{ID: 12, FirstName: "Max", LastName: "Test", LanguageCode: "code1"}, "Max Test", 12},
//		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "nickName", LanguageCode: "code1"}, "nickName", 13},
//		{&tgbotapi.User{}, " ", 0},
//	}
//
//	for _, pair := range tests {
//		resultID, resultStr := GetNickName(pair.input)
//		if resultStr != pair.outputStr || resultID != pair.outputInt {
//			t.Errorf("For %v\nexpected %s %d\ngot %s %d", pair.input, pair.outputStr, pair.outputInt, resultStr, resultID)
//		}
//	}
//}
