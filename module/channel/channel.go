package channel

import (
	"sync"
)

var channel chan []byte

//用来通知channel 写完的channel
var oc chan []byte
var channelManager sync.Once
var ocManager sync.Once

func GetChannel() chan []byte {
	channelManager.Do(func() {
		channel = make(chan []byte, 100)
	})
	return channel
}

func GetOverChannel() chan []byte {
	ocManager.Do(func() {
		oc = make(chan []byte, 100)
	})
	return oc
}
