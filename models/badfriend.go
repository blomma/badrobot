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

func (r *BadFriends) SetResult(value []byte) {
	r.result.Lock()
	defer r.result.Unlock()
	r.result.value = value
}

func (r *BadFriends) Result() []byte {
	r.result.RLock()
	defer r.result.RUnlock()
	return r.result.value
}

type result struct {
	sync.RWMutex
	value []byte
}

type BadFriends struct {
	result *result
}

func NewBadFriends() *BadFriends {
	bf := &BadFriends{
		result: &result{},
	}

	go fetchBadFriends(bf)

	return bf
}

type badFriend struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func init() {
	localPool, err := pool.New("tcp", "192.168.1.5:6379", 10)
	if err != nil {
		log.Panic(err)
	}

	globalPool = localPool
}

func fetchBadFriends(b *BadFriends) {
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

			b.SetResult(jsonBadFriends)
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
