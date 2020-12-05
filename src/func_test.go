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
func TestGetListHelps(t *testing.T) {
	t.Parallel()

	type testPair struct {
		telegramNickName *tgbotapi.User
		ownNickName      string
		output           string
	}

	var tests = []testPair{
		{&tgbotapi.User{ID: 12, FirstName: "Max", LastName: "Test", LanguageCode: "code1"}, "own1", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n"},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "nickName", LanguageCode: "code1"}, "own1", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n"},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "own1", LanguageCode: "code1"}, "own1", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n/show - показать все коды;\n/reset - удалить все из БД и создать новые;\n/add - добавить новые правильные коды в формате: Code,Danger,Sector;\n/update - обновить коды в бд, в формате: CodeNew,Danger,Sector,CodeOld;\n/delete - удалить указанный код;\n"},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "own1", LanguageCode: "code1"}, "", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n"},
		{&tgbotapi.User{}, "own1", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n"},
		{&tgbotapi.User{}, "", "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n/show - показать все коды;\n/reset - удалить все из БД и создать новые;\n/add - добавить новые правильные коды в формате: Code,Danger,Sector;\n/update - обновить коды в бд, в формате: CodeNew,Danger,Sector,CodeOld;\n/delete - удалить указанный код;\n"},
	}

	for _, pair := range tests {
		result := GetListHelps(pair.telegramNickName, pair.ownNickName)
		if result != pair.output {
			t.Errorf("For %v\nexpected %s\ngot %s", pair.telegramNickName, pair.output, result)
		}
	}
}
