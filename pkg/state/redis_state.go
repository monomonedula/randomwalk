package state

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisState struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisState(
	rdb *redis.Client,
	ctx context.Context,
) RedisState {
	return RedisState{rdb, ctx}
}

func (rs *RedisState) UserData(userId int) (UserDataStruct, error) {
	userData, err := rs.rdb.HGetAll(rs.ctx, strconv.Itoa(userId)).Result()
	if err != nil {
		return UserDataStruct{}, nil
	}

	if dataStr, ok := userData["data"]; ok {
		data := UserDataStruct{}
		err := json.Unmarshal([]byte(dataStr), &data)
		return data, err
	}
	return UserDataStruct{}, nil
}

func (rs *RedisState) UpdateUserData(userID int, newData UserDataStruct) error {
	j, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	err = rs.rdb.HSet(rs.ctx, strconv.Itoa(userID), "data", string(j)).Err()
	if err != nil {
		return err
	}
	return nil
}
