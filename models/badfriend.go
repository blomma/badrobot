package models

import (
	"encoding/json"
	"log"

	"github.com/mediocregopher/radix.v2/pool"
)

type BadFriend struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var g_db *pool.Pool

func init() {
	db, err := pool.New("tcp", "192.168.1.5:6379", 10)
	if err != nil {
		log.Panic(err)
	}

	g_db = db
}

func GetAllBadFriends() ([]*BadFriend, error) {
	conn, err := g_db.Get()
	if err != nil {
		return nil, err
	}
	defer g_db.Put(conn)

	ids, err := conn.Cmd("SMEMBERS", "badfriends").List()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(ids); i++ {
		ids[i] = "badfriend:" + ids[i]
	}

	replyBadfriends, err := conn.Cmd("MGET", ids).ListBytes()
	if err != nil {
		return nil, err
	}

	badfriends := make([]*BadFriend, len(ids))
	for i, value := range replyBadfriends {
		badfriend := new(BadFriend)
		if err := json.Unmarshal(value, badfriend); err != nil {
			return nil, err
		}

		badfriends[i] = badfriend
	}

	return badfriends, nil
}
