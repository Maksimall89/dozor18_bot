package src

import (
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telegram-bot-api.v4"
	"testing"
)

const errorExpect = "For %v\nexpected %s\ngot %s"

func TestGetNickName(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input     *tgbotapi.User
		outputStr string
		outputInt int
	}

	var tests = []testPair{
		{&tgbotapi.User{ID: 12, FirstName: "Max", LastName: "Test", LanguageCode: "code1"}, "Max Test", 12},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "nickName", LanguageCode: "code1"}, "nickName", 13},
		{&tgbotapi.User{}, " ", 0},
	}

	for _, pair := range tests {
		resultID, resultStr := GetNickName(pair.input)
		if resultStr != pair.outputStr || resultID != pair.outputInt {
			t.Errorf("For %v\nexpected %s %d\ngot %s %d", pair.input, pair.outputStr, pair.outputInt, resultStr, resultID)
		}
	}
}
func TestGetMD5Hash(t *testing.T) {
	t.Parallel()

	input := "test"
	output := "098f6bcd4621d373ca"
	result := GetMD5Hash(input)
	if result != output {
		t.Errorf(errorExpect, input, output, result)
	}
}
func TestGetListHelps(t *testing.T) {
	t.Parallel()

	type testPair struct {
		telegramNickName *tgbotapi.User
		OwnID            int
		output           string
	}

	userHelps := "/help - информация по всем доступным командам;\n/codes - коды;\n/generate, /gen - сгенерировать коды;\n/text - текст приквела;\n/create - создать команду;\n/join - вступить в команду;\n/list - список участников команды;\n/listusers - список участников в командах;\n/leave - выйти из команды;\n/invite - получить ссылку приглашение в команду;\n"
	adminHelps := userHelps + "/show - показать все коды;\n/reset - удалить все из БД и создать новые;\n/add - добавить новые правильные коды в формате: Code,Danger,Sector;\n/update - обновить коды в бд, в формате: CodeNew,Danger,Sector,CodeOld;\n/delete - удалить указанный код;\n/listteams - список всех команд;\n/resetteams - удалить все команды;\n"

	var tests = []testPair{
		{&tgbotapi.User{ID: 12, FirstName: "Max", LastName: "Test", LanguageCode: "code1"}, 13, userHelps},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "nickName", LanguageCode: "code1"}, 11, userHelps},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "own1", LanguageCode: "code1"}, 13, adminHelps},
		{&tgbotapi.User{ID: 13, FirstName: "Max", LastName: "Test", UserName: "own1", LanguageCode: "code1"}, 0, userHelps},
		{&tgbotapi.User{}, 13, userHelps},
		{&tgbotapi.User{}, 0, adminHelps},
	}

	for _, pair := range tests {
		result := GetListHelps(pair.telegramNickName, pair.OwnID)
		if result != pair.output {
			t.Errorf(errorExpect, pair.telegramNickName, pair.output, result)
		}
	}
}
func TestCheckMessage(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  string
		output string
	}

	rightErr := ("&#10071;Недопустимые символы в сообщении. Можно использовать лишь буквы и цифры русского и английского алфавита.")
	var tests = []testPair{
		{"dfg=dfg", rightErr},
		{"5*5=10", rightErr},
	}
	for _, pair := range tests {
		result := CheckMessage(pair.input)
		assert.EqualErrorf(t, result, rightErr, errorExpect, pair.input, pair.output, result)
	}

	tests = []testPair{
		{"qwe", ""},
		{"qwfsdfSSFS efsdf", ""},
	}
	for _, pair := range tests {
		result := CheckMessage(pair.input)
		if result != nil {
			t.Errorf(errorExpect, pair.input, pair.output, result)
		}
	}
}
