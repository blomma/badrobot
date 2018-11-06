package models

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

var stopChan = make(chan struct{})
var stoppedChan = make(chan struct{})

func (r *BadFriendsModel) SetResult(value []byte) {
	r.result.Lock()
	defer r.result.Unlock()
	r.result.value = value
}

func (r *BadFriendsModel) Result() []byte {
	r.result.RLock()
	defer r.result.RUnlock()
	return r.result.value
}

type result struct {
	sync.RWMutex
	value []byte
}

type BadFriendsModel struct {
	result *result
}

func NewBadFriendsModel(redisServer string) (*BadFriendsModel, func()) {
	badFriends := &BadFriendsModel{
		result: &result{},
	}

	go fetchBadFriends(badFriends, redisServer)

	stop := func() {
		close(stopChan)
		<-stoppedChan
	}

	return badFriends, stop
}

type badFriend struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func fetchBadFriends(b *BadFriendsModel, redisServer string) {
	defer close(stoppedChan)
	for {
		badFriends, err := getAllBadFriends(redisServer)
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

		timerChan := make(chan struct{})
		go func() {
			<-time.After(5 * time.Minute)
			close(timerChan)
		}()

		select {
		case <-stopChan:
			return
		case <-timerChan:
			break
		}
	}
}

func getAllBadFriends(redisServer string) ([]*badFriend, error) {
	client, err := redis.Dial("tcp", redisServer)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ids, err := client.Cmd("SMEMBERS", "badfriends").List()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(ids); i++ {
		ids[i] = "badfriend:" + ids[i]
	}

	replyBadfriends, err := client.Cmd("MGET", ids).ListBytes()
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
