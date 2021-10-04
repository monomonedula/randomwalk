package texter

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TexterSimple struct {
	ttable      map[string]map[string]string
	defaultLang string
}

func (t *TexterSimple) TextFor(user tb.User, code string) string {
	langTable, ok := t.ttable[user.LanguageCode]
	if !ok {
		langTable = t.ttable[t.defaultLang]
	}
	if reply, ok := langTable[code]; ok {
		return reply
	}
	panic(fmt.Sprintf("No text for %s code %s", user.LanguageCode, code))
}

func BasicTexter() *TexterSimple {
	return &TexterSimple{
		defaultLang: "en",
		ttable: map[string]map[string]string{
			"en": {
				"whoops": "Whoops. Something went wrong. Please try again from the start :(",
				"help": "I generate random point on a map for you to walk to. " +
					"Firstly, send me your current location, then text me min and max distance so I could estimate the range" +
					" you want the destination to be located within. \n" +
					"Note that the distances should specified in meters as integers (e.g. 4000, 2500, 500, etc.)",

				"bad_range_format":        "Looks like the format was wrong. It should be two numbers, something like this: 500 2300",
				"input_range":             "OK. Now send me the min and max distance you'd like to walk in meters separated by a space (Like this: 1000 2500)",
				"send_you_location_first": "Send your location first",
			},
			"ru": {
				"whoops": "Ой, что-то пошло не так. Попробуйте еще раз сначала :(",
				"help": "Я генерирую случайные точки для прогулок. " +
					"Сначала вышли мне своё местоположение, затем укажи минимальное и макисмальное расстояние чтобы я" +
					" понял на какое расстояние прогулки ты рассчитываешь. \n" +
					"Важно: расстояние следует указывать в метрах как целые числа (например: 4000, 2500, 500 и т.п.)",

				"bad_range_format":        "Похоже на ошибку. Укажи минимальное и максимальное расстояние через пробел примерно вот так: 3000 5000",
				"input_range":             "Ок. Теперь отправь мне минимальное и максимальное расстояние до пункта назначения в метрах через пробел (вроде этого: 1000 2500)",
				"send_you_location_first": "Сперва вышли своё местоположение",
			},
			"uk": {
				"whoops": "Ой, щось пішло не так. Спробуйте ще раз спочатку :(",
				"help": "Я генерую випадкові місця для прогулянок. " +
					"Спочатку вишли мені своє місцезнаходження, потім вкажи мінімальну і максимальну відстань до пункту призначення. " +
					"Важливо: відстань слід вказати у метрах через пробіл (ось так: 1000 2500)",

				"bad_range_format":        "Схоже на помилку. Вкажи мінімальну і максимальну відстань через пробіл приблизно ось так: 3000 5000",
				"input_range":             "Ок. Тепер відправ мені мінімальну та максимальну відстань до пункту призначення у метрах через пробіл (наприклад, отак: 1000 2500)",
				"send_you_location_first": "Спершу надішли своє місцезнаходження",
			},
		},
	}
}
