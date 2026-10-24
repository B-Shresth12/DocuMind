package response

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Data(ctx *gin.Context, status int, message string, data any) {
	ctx.Header("Content-Type", "application/json")

	ctx.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

//in case only message needs to be send especially for sending errors
type ResponseWithMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Message(ctx *gin.Context, status int, message string) {
	ctx.Header("Content-Type", "application/json")

	ctx.JSON(status, ResponseWithMessage{
		Status:  status,
		Message: message,
	})
}
