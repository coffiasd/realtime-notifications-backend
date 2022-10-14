package events

import (
	"context"
	"encoding/json"
	"fmt"
	"notify-server/middleware"
)

var ctx = context.Background()

const (
	SenderTransactionCompleteEvent   = 1
	ReceiverTransactionCompleteEvent = 2
	EVENT_TOPIC                      = "event_topic"
)

type Subscription struct {
	Address  string `json:"address"`
	DeviceId string `json:"device_id"`
	Email    string `json:"email"`
	EventId  int    `json:"event-id"`
}

type FetchResponse struct {
	CurrentRound int           `json:"current-round"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id             string            `json:"id"`
	TxType         string            `json:"tx-type"`
	ConfirmedRound int               `json:"confirmed-round"`
	Sender         string            `json:"sender"`
	Fee            int               `json:"fee"`
	RoundTime      int               `json:"round-time"`
	Payment        PaymentTransation `json:"payment-transaction"`
}

type PaymentTransation struct {
	Amount   int    `json:"amount"`
	Receiver string `json:"receiver"`
}

func GetSenderRet(address string) (*Subscription, error) {
	var sub Subscription
	sender, err := middleware.GetRedis().HGet(ctx, GetTopicByEventId(SenderTransactionCompleteEvent), address).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(sender), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

func GetReceiverRet(address string) (*Subscription, error) {
	var sub Subscription
	sender, err := middleware.GetRedis().HGet(ctx, GetTopicByEventId(ReceiverTransactionCompleteEvent), address).Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(sender), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

func GetTopicByEventId(id int) string {
	switch id {
	case SenderTransactionCompleteEvent:
		return fmt.Sprintf("topic:%d", SenderTransactionCompleteEvent)
	case ReceiverTransactionCompleteEvent:
		return fmt.Sprintf("topic:%d", ReceiverTransactionCompleteEvent)
	default:
		return "unknow"
	}
}
