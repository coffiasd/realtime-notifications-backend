package cron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"notify-server/config"
	"notify-server/events"
	"notify-server/middleware"
	"notify-server/utils"
	"time"

	"github.com/gregdel/pushover"
)

type HealthJson struct {
	Round int `json:"round"`
}

func health() (int, error) {
	var ret HealthJson

	//get current round
	healthJson, err := utils.Request(config.ParseConfig.App.AlgoRandUrl+"/health", time.Second*5)
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal(healthJson, &ret); err != nil {
		return 0, errors.New("get health json Unmarshal error")
	}

	return ret.Round, nil
}

func Fetch() error {
	var trans events.FetchResponse
	round, err := health()
	if err != nil || round == 0 {
		return err
	}

	for {
		trans = events.FetchResponse{}
		FetchJson, err := utils.Request(fmt.Sprintf(config.ParseConfig.App.AlgoRandUrl+"/v2/transactions?round=%d", round), time.Second*5)

		middleware.GetLogger().Info(time.Now().Unix(), ",round:", round)

		if err != nil {
			//damn idk what's happening.
			middleware.GetLogger().Error(time.Now().Unix(), "Request error:", err)
			continue
		}

		if err := json.Unmarshal(FetchJson, &trans); err != nil {
			middleware.GetLogger().Error(time.Now().Unix(), "json Unmarshal:", err)
			continue
		}

		if trans.CurrentRound == round {
			middleware.GetLogger().Info(time.Now().Unix(), ",sleeping")
			time.Sleep(time.Second * 2)
			continue
		}

		// fmt.Println(trans.Transactions)
		for _, v := range trans.Transactions {
			// middleware.GetLogger().Info("INFO:", v)
			checkPushToList(v)
		}

		if round < trans.CurrentRound {
			round++
		}

		time.Sleep(time.Second * 2)
	}
}

func checkPushToList(tran events.Transaction) {
	var eventIds []int
	senderData, _ := middleware.GetRedis().HGet(context.Background(), events.EVENT_TOPIC, tran.Sender).Result()
	if senderData != "" {
		json.Unmarshal([]byte(senderData), &eventIds)
		for _, v := range eventIds {
			if v == events.ReceiverTransactionCompleteEvent {
				//push to notify list
				// fmt.Println("sender:", v, events.ReceiverTransactionCompleteEvent)
				middleware.GetLogger().Info(time.Now().Unix(), ",pushing senderData", v)
				go pushNotify(tran.Sender, getSendContent(tran))
			}
		}
	}

	receiverData, _ := middleware.GetRedis().HGet(context.Background(), events.EVENT_TOPIC, tran.Payment.Receiver).Result()
	if receiverData != "" {
		json.Unmarshal([]byte(receiverData), &eventIds)
		for _, v := range eventIds {
			if v == events.SenderTransactionCompleteEvent {
				//push to notify list
				// fmt.Println("receive:", v, events.ReceiverTransactionCompleteEvent)
				middleware.GetLogger().Info(time.Now().Unix(), ",pushing receiverData", v)
				go pushNotify(tran.Payment.Receiver, getReceiveContent(tran))
			}
		}
	}

}

func pushNotify(address string, content string) {
	// Create a new pushover app with a token
	app := pushover.New("aqu43exja3ncvm7xt9q43gc593y2hh")
	//get push ID.
	pushId, _ := middleware.GetRedis().Get(context.Background(), address).Result()
	// Create a new recipient
	recipient := pushover.NewRecipient(pushId)
	// Create the message to send
	message := pushover.NewMessage(content)
	// Send the message to the recipient
	app.SendMessage(message, recipient)
}

//sender format
func getSendContent(tran events.Transaction) string {
	completeDate := time.Unix(int64(tran.RoundTime), 0).Format("2006/01/02 15:04:05")
	return fmt.Sprintf("Sender: your transaction tx-type:%s,confirmed-round:%d,fee:%d is complete at %s", tran.TxType, tran.ConfirmedRound, tran.Fee, completeDate)
}

//receiver format
func getReceiveContent(tran events.Transaction) string {
	completeDate := time.Unix(int64(tran.RoundTime), 0).Format("2006/01/02 15:04:05")
	return fmt.Sprintf("Receiver: your transaction tx-type:%s,confirmed-round:%d,fee:%d is complete at %s ", tran.TxType, tran.ConfirmedRound, tran.Fee, completeDate)
}
