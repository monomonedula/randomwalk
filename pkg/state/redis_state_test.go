package state

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tucnak/telebot.v2"
)

var ctx = context.TODO()

func TestRedisStateUserDataNewUser(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userId := 1234532920
	rs := NewRedisState(db, ctx)

	mock.ExpectHGetAll(strconv.Itoa(userId))
	data, err := rs.UserData(userId)

	assert.Nil(t, err)
	assert.Equal(t, UserDataStruct{}, data)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRedisStateUserExisting(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userId := 1234532920
	rs := NewRedisState(db, ctx)

	expected := mock.ExpectHGetAll(strconv.Itoa(userId))
	expected.SetVal(
		map[string]string{
			"data": `{"Location": {"latitude": 51.491524, "longitude": -0.1278841}}`,
		},
	)

	data, err := rs.UserData(userId)

	assert.Nil(t, err)
	assert.Equal(t, UserDataStruct{Location: telebot.Location{Lat: 51.491524, Lng: -0.1278841}}, data)

	assert.Nil(t, mock.ExpectationsWereMet())

}

func TestRedisStateUserUpdate(t *testing.T) {
	db, mock := redismock.NewClientMock()
	userId := 1234532920
	rs := NewRedisState(db, ctx)

	exp := mock.ExpectHSet(strconv.Itoa(userId), "data", `{"Location":{"latitude":51.491524,"longitude":-0.1278841}}`)
	exp.SetVal(0)

	err := rs.UpdateUserData(userId, UserDataStruct{Location: telebot.Location{Lat: 51.491524, Lng: -0.1278841}})

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())

}
