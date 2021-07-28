package eventer

import (
	"github.com/hoank92/hlibs/client"
	"github.com/hoank92/hlibs/network/middleware"
)

type Client struct {
	client.Request
}

type Status struct {
	Status bool `json:"status"`
}

type Event struct {
	Name      string            `json:"name"`
	CreatedAt int64             `json:"created_at"`
	Params    map[string]string `json:"params"`
}

var senderClient *Client
var batchLimit = 100

func Init(url string, limit int) {
	senderClient = &Client{client.InitJSONClient(url, middleware.SentryErrorHandler)}
	batchLimit = limit
}

func (c Client) SendBatch(events []Event) (status Status, err error) {
	err = senderClient.Post(&status, "", events)
	return
}
