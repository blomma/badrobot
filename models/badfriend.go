package models

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

const redisHost string = "192.168.1.5:6379"

var stopChan = make(chan struct{})
var stoppedChan = make(chan struct{})

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

func NewBadFriends() (*BadFriends, func()) {
	b := &BadFriends{
		result: &result{},
	}

	go fetchBadFriends(b)

	stop := func() {
		close(stopChan)
		<-stoppedChan
	}

	return b, stop
}

type badFriend struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func fetchBadFriends(b *BadFriends) {
	defer close(stoppedChan)
	for {
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

func getAllBadFriends() ([]*badFriend, error) {
	client, err := redis.Dial("tcp", redisHost)
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
