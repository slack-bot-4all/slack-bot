package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/slack-bot-4all/slack-bot/src/model"
	"github.com/slack-bot-4all/slack-bot/src/service"
)

// AddRancher : add a new Rancher to db
func AddRancher(c *gin.Context) {
	var r model.Rancher
	c.BindJSON(&r)

	err := service.AddRancher(&r)
	if err != nil {
		ResponseJSON(c, 400, nil)
	} else {
		ResponseJSON(c, 200, r)
	}
}

// ListRancher : list all ranchers
func ListRancher(c *gin.Context) {
	ranchers, err := service.ListRancher()

	if err != nil {
		ResponseJSON(c, 404, nil)
	} else {
		ResponseJSON(c, 200, ranchers)
	}
}
