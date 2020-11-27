package src

import (
	"gopkg.in/telegram-bot-api.v4"
	"testing"
)

func TestGetNickName(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  *tgbotapi.User
		output string
	}

	var tests = []testPair{
		{&tgbotapi.User{ID: 12, FirstName: "Max", LastName: "Test", LanguageCode: "code1"}, "12 + Max Test"},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "nickName", LanguageCode: "code1"}, "13 + nickName"},
		{&tgbotapi.User{}, "0 +  "},
	}

	for _, pair := range tests {
		result := GetNickName(pair.input)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
