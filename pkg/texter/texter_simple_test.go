package texter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
)

func TestBasicTexterKeysPresent(t *testing.T) {
	txtr := BasicTexter()

	for _, v := range []string{"whoops", "help", "bad_range_format", "input_range"} {
		for _, lang := range []string{"en", "ru", "uk"} {
			s := txtr.TextFor(telebot.User{LanguageCode: lang}, v)
			assert.NotEqual(t, "", s)
		}
	}

	assert.Equal(t, "Whoops. Something went wrong. Please try again from the start :(", txtr.TextFor(telebot.User{LanguageCode: "en"}, "whoops"))
	assert.Equal(t, "Ой, что-то пошло не так. Попробуйте еще раз сначала :(", txtr.TextFor(telebot.User{LanguageCode: "ru"}, "whoops"))
	assert.Equal(t, "Ой, щось пішло не так. Спробуйте ще раз спочатку :(", txtr.TextFor(telebot.User{LanguageCode: "uk"}, "whoops"))
}

func TestBasicTexterDefaultLang(t *testing.T) {
	txtr := BasicTexter()

	for _, v := range []string{"whoops", "help", "bad_range_format", "input_range"} {
		assert.Equal(t, txtr.TextFor(telebot.User{LanguageCode: "en"}, v), txtr.TextFor(telebot.User{LanguageCode: "whatever"}, v))
	}
}
