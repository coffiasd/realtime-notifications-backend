package api

import (
	"context"
	"encoding/json"
	"errors"
	"notify-server/events"
	"notify-server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var ctx = context.Background()

type postSubParams struct {
	Address  string `json:"address" validate:"required,max=64"`
	DeviceId string `json:"device_id" validate:"required,max=64"`
}

type saveSubParams struct {
	Address  string `json:"address" validate:"required,max=64"`
	EventIds []int  `json:"event_id" validate:"required"`
}

// AddSub
// @address
// @deviceId
// @email
func SaveInfo(c *gin.Context) {
	var params postSubParams

	//json marshal error
	if err := c.ShouldBindJSON(&params); err != nil {
		JsonReturn(c, "params error", err)
		return
	}

	//validator params
	validator := validator.New()
	if err := validator.Struct(&params); err != nil {
		JsonReturn(c, "params error", err)
		return
	}

	middleware.GetRedis().Set(ctx, params.Address, params.DeviceId, 0)

	JsonReturn(c, "", nil)
}

// UserInfo
// @address
func UserInfo(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		JsonReturn(c, "", errors.New("address required"))
		return
	}

	info, _ := middleware.GetRedis().Get(ctx, address).Result()

	JsonReturn(c, info, nil)
}

// SubList
// @address
func SubList(c *gin.Context) {
	var list []int
	address := c.Query("address")
	if address == "" {
		JsonReturn(c, "", errors.New("address required"))
		return
	}

	data, _ := middleware.GetRedis().HGet(ctx, events.EVENT_TOPIC, address).Result()
	json.Unmarshal([]byte(data), &list)

	JsonReturn(c, list, nil)
}

// DelSub
// @address
// @eventId
func SaveSub(c *gin.Context) {
	var params saveSubParams

	//json marshal error
	if err := c.ShouldBindJSON(&params); err != nil {
		JsonReturn(c, "params error", err)
		return
	}

	//validator params
	validator := validator.New()
	if err := validator.Struct(&params); err != nil {
		JsonReturn(c, "params error", err)
		return
	}

	//check user info is saved?
	user, _ := middleware.GetRedis().Get(ctx, params.Address).Result()
	if user == "" {
		JsonReturn(c, "user info error", errors.New("user info error"))
		return
	}

	jsonData, _ := json.Marshal(params.EventIds)
	middleware.GetRedis().HSet(ctx, events.EVENT_TOPIC, params.Address, string(jsonData))

	JsonReturn(c, "", nil)
}
