package utils

import (
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ReturnFormattedData(ctx *gin.Context, code int, msg string, data interface{}) {
	res := &ResponseBody{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	ctx.AsciiJSON(code, res)
}
