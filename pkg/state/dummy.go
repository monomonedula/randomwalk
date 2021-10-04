package state

type dummyState struct {
	Data map[int]UserDataStruct
}

func (ds *dummyState) UserData(userId int) (UserDataStruct, error) {
	if data, ok := ds.Data[userId]; ok {
		return data, nil
	}
	return UserDataStruct{}, nil
}

func (ds *dummyState) UpdateUserData(userID int, newData UserDataStruct) error {
	ds.Data[userID] = newData
	return nil
}

func NewDummyState() *dummyState {
	return &dummyState{make(map[int]UserDataStruct)}
}
