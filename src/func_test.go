package src

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
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
	key := "myKey"
	output := "0d23530616c7d05f6e"
	result := GetMD5Hash(input, key)
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

	userHelps := "/help - все команды;\n" +
		"/codes - список кодов;\n" +
		"/gen - сгенерировать коды;\n" +
		"/text - текст задания;\n" +
		"/create имя команды - создать команду;\n" +
		"/join - вступить в команду;\n" +
		"/list - участники команды;\n" +
		"/listusers - участники в командах;\n" +
		"/leave - выйти из команды;\n" +
		"/invite - приглашение в команду;\n" +
		"/teams - список команд;\n"
	adminHelps := userHelps +
		"/show - показать коды;\n" +
		"/reset - удалить данные из <b>teams</b> или <b>codes</b>;\n" +
		"/add - добавить код: <b>Code,Danger,Sector,TimeBonus,Tasks</b>;\n" +
		"/update - обновить коды: <b>CodeNew,Danger,Sector,TimeBonus,TaskID,CodeOld</b>;\n" +
		"/delete - удалить код;\n" +
		"/listteams - список команд;\n" +
		"/createdb - создать в БД;\n" +
		"/createtask - создать задание;\n"

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

	errSymbol := "&#10071;Недопустимые символы в сообщении. Можно использовать лишь буквы и цифры русского и английского алфавита"
	errMinLen := "&#10071;Сообщение слишком короткое"
	errMaxLen := "&#10071;Сообщение слишком длинное"
	var tests = []testPair{
		{"dfg=dfg", errSymbol},
		{"5*5=10", errSymbol},
		{"a", errMinLen},
		{"bb", errMinLen},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", errMaxLen},
	}
	for _, pair := range tests {
		result := CheckMessage(pair.input)
		assert.EqualErrorf(t, result, pair.output, errorExpect, pair.input, pair.output, result)
	}

	tests = []testPair{
		{"qwe", ""},
		{"qwfsdfSSFS efsdf", ""},
		{"qwf123gdf5g", ""},
		{"qwfhr63", ""},
		{"код14444", ""},
	}
	for _, pair := range tests {
		result := CheckMessage(pair.input)
		if result != nil {
			t.Errorf(errorExpect, pair.input, pair.output, result)
		}
	}
}
func TestArrTrimSpace(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  []string
		output []string
	}

	rightArr := []string{"a", "b", "c", "d"}
	var tests = []testPair{
		{[]string{" a", "b ", " c ", "d"}, rightArr},
		{[]string{"a", " b b b", "c c ", "d e"}, []string{"a", "b b b", "c c", "d e"}},
		{[]string{"a", "b", "c", "d"}, rightArr},
	}

	for _, pair := range tests {
		ArrTrimSpace(pair.input)
		for number, value := range pair.input {
			if pair.output[number] != value {
				t.Errorf("For %v\nexpected %s", pair.output[number], value)
			}
		}
	}
}
func TestGeneralConvertTimeSec(t *testing.T) {
	t.Parallel()

	type testPair struct {
		input  int
		output string
	}

	var tests = []testPair{
		{0, "0 секунд"},
		{1, "1 секунда"},
		{60, "1 минута"},
		{66, "1 минута 6 секунд"},
		{120, "2 минуты"},
		{122, "2 минуты 2 секунды"},
		{600, "10 минут"},
		{3600, "1 час"},
		{3601, "1 час 1 секунда"},
		{3661, "1 час 1 минута 1 секунда"},
		{86400, "1 день"},
		{90061, "1 день 1 час 1 минута 1 секунда"},
		{36045645, "417 дней 4 часа 45 минут 45 секунд"},
	}

	for _, pair := range tests {
		result := ConvertTimeSec(pair.input)
		if result != pair.output {
			t.Errorf("For %d\nexpected %s\ngot %s", pair.input, pair.output, result)
		}
	}
}
