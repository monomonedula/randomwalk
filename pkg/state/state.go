package state

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type State interface {
	UserData(userId int) (UserDataStruct, error)
	UpdateUserData(userID int, newData UserDataStruct) error
}

type UserDataStruct struct {
	Location tb.Location
}
