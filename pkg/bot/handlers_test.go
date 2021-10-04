package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/monomonedula/randomwalk/pkg/state"
	"github.com/monomonedula/randomwalk/pkg/texter"
	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

type ComparableCall struct {
	To   tb.Recipient
	What interface{}
}

type SendCall struct {
	to      tb.Recipient
	what    interface{}
	options []interface{}
}

type SendCallExpected struct {
	SendCall

	msg *tb.Message
	err error
}

type mock struct {
	expected []SendCallExpected
	actual   []SendCall
}

func (m *mock) Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	m.actual = append(m.actual, SendCall{to, what, options})
	current := ComparableCall{to, what}
	for _, x := range m.expected {
		sc := ComparableCall{x.to, x.what}
		if sc == current {
			return x.msg, x.err
		}
	}
	return &tb.Message{}, nil
}

func (m *mock) ExpectSend(call SendCallExpected) {
	m.expected = append(m.expected, call)
}

func (m *mock) ExpectationsWereMet() error {
	if len(m.expected) != len(m.actual) {
		return errors.New(fmt.Sprintf("Expected %d, got %d calls", len(m.expected), len(m.actual)))
	}
	for i := 0; i < len(m.expected); i++ {
		x := m.expected[i]
		ec := ComparableCall{x.to, x.what}
		ac := ComparableCall{m.actual[i].to, m.actual[i].what}
		expectedStr, _ := json.Marshal(ec)
		actualStr, _ := json.Marshal(ac)
		if string(expectedStr) != string(actualStr) {
			return errors.New(fmt.Sprintf("Expected %s got %s", string(expectedStr), string(actualStr)))
		}
	}

	return nil
}

func TestHelp(t *testing.T) {
	s := mock{}
	usr := &tb.User{LanguageCode: "ru", ID: 24131356}
	txt := texter.BasicTexter().TextFor(*usr, "help")
	s.ExpectSend(
		SendCallExpected{SendCall: SendCall{to: usr, what: txt}},
	)
	Help(&s, texter.BasicTexter())(&tb.Message{Sender: usr})
	assert.Nil(t, s.ExpectationsWereMet())
}

func TestAcceptTextNoLocation(t *testing.T) {
	s := mock{}
	usr := &tb.User{LanguageCode: "ru", ID: 24131356}
	st := state.NewDummyState()
	s.ExpectSend(
		SendCallExpected{SendCall: SendCall{to: usr, what: "Сперва вышли своё местоположение"}},
	)
	AcceptText(&s, st, texter.BasicTexter(), ErrorHandler(&s, texter.BasicTexter(), *log.Default()))(
		&tb.Message{Sender: usr, Text: "300 700"},
	)
	assert.Nil(t, s.ExpectationsWereMet())
}

func TestAcceptTextOk(t *testing.T) {
	s := mock{}
	usr := &tb.User{LanguageCode: "ru", ID: 24131356}
	st := state.NewDummyState()
	st.UpdateUserData(usr.ID, state.UserDataStruct{Location: tb.Location{Lat: 55.75631065468432, Lng: 37.608801439307925}})
	s.ExpectSend(
		SendCallExpected{SendCall: SendCall{to: usr, what: &tb.Location{Lat: 55.759495, Lng: 37.606586}}},
	)
	AcceptText(&s, st, texter.BasicTexter(), ErrorHandler(&s, texter.BasicTexter(), *log.Default()))(
		&tb.Message{Sender: usr, Text: "300 700"},
	)
	assert.Nil(t, s.ExpectationsWereMet())
}

func TestAcceptLocation(t *testing.T) {
	s := mock{}
	usr := &tb.User{LanguageCode: "ru", ID: 24131356}
	st := state.NewDummyState()
	txt := "Ок. Теперь отправь мне минимальное и максимальное расстояние до пункта назначения в метрах через пробел (вроде этого: 1000 2500)"
	s.ExpectSend(
		SendCallExpected{SendCall: SendCall{to: usr, what: txt}},
	)
	AcceptLocation(&s, st, texter.BasicTexter(), ErrorHandler(&s, texter.BasicTexter(), *log.Default()))(
		&tb.Message{Sender: usr, Location: &tb.Location{Lat: 55.75631065468432, Lng: 37.608801439307925}},
	)
	assert.Nil(t, s.ExpectationsWereMet())
	data, err := st.UserData(usr.ID)
	assert.Nil(t, err)
	assert.Equal(t, state.UserDataStruct{Location: tb.Location{Lat: 55.75631065468432, Lng: 37.608801439307925}}, data)
}
