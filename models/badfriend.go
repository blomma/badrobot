package models

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
)

var globalPool *pool.Pool
var globalStopChan = make(chan struct{})

var BadFriends = &badFriendsResult{}

type badFriend struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type badFriendsResult struct {
	sync.RWMutex
	value []byte
}

func (b *badFriendsResult) Set(value []byte) {
	b.Lock()
	defer b.Unlock()
	b.value = value
}

func (b *badFriendsResult) Get() []byte {
	b.RLock()
	defer b.RUnlock()
	return b.value
}

func init() {
	localPool, err := pool.New("tcp", "192.168.1.5:6379", 10)
	if err != nil {
		log.Panic(err)
	}

	globalPool = localPool
	go fetchBadFriends(BadFriends)
}

func fetchBadFriends(b *badFriendsResult) {
	for {
		select {
		default:
			badFriends, err := getAllBadFriends()
			if err != nil {
				log.Println(err)
				break
			}

			jsonBadFriends, err := json.Marshal(badFriends)
			if err != nil {
				log.Println(err)
				break
			}

			b.Set(jsonBadFriends)
		case <-globalStopChan:
			return
		}

		time.Sleep(5 * time.Minute)
	}
}

func getAllBadFriends() ([]*badFriend, error) {
	conn, err := globalPool.Get()
	if err != nil {
		return nil, err
	}
	defer globalPool.Put(conn)

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

	badfriends := make([]*badFriend, len(ids))
	for i, value := range replyBadfriends {
		badfriend := new(badFriend)
		if err := json.Unmarshal(value, badfriend); err != nil {
			return nil, err
		}

		badfriends[i] = badfriend
	}

	return badfriends, nil
}
