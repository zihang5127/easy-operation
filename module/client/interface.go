package client

import "net/url"

type Interface interface {
	Command(host url.URL, username, password, shell string, channel chan<- []byte)
}
