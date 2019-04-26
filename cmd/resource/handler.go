package resource

import (
	"github.com/gin-gonic/gin"
)

// Response : Struct default to responses
type Response struct {
	Status  int         `json:"statusCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseJSON : function to returns a response
func ResponseJSON(w *gin.Context, status int, payload interface{}) {
	var resp Response
	resp.Status = status
	resp.Data = payload

	switch status {
	case 200:
		resp.Message = "Successful request"
	case 400:
		resp.Message = "Bad request"
	case 404:
		resp.Message = "Resource not found"
	}

	w.JSON(status, resp)
}
