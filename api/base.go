package api

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse() *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
	}
}

func resolveError(err error) (int, string) {
	if err == nil {
		return 200, "success"
	}

	return 500, err.Error()
}

func JsonReturn(c *gin.Context, d interface{}, err error) {
	res := NewResponse()
	if err != nil {
		res.Code, res.Msg = resolveError(err)
	}
	res.Data = d
	c.JSON(200, res)
}
