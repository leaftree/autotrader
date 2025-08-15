package wxpusher

import (
	"fmt"

	"github.com/leaftree/autotrader/config"
	log "github.com/leaftree/autotrader/logger"
	wxpusher "github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

var (
	logger = log.NewLogger("wxpusher")
)

func SendNotify(msg string) {
	appToken := config.GetConfig().WxPusher.AppToken
	userID := config.GetConfig().WxPusher.UserID

	content := model.NewMessage(appToken).SetContent(msg).AddUId(userID)
	resp, err := wxpusher.SendMessage(content)
	if err != nil {
		fmt.Println(resp, err)
	}
	logger.Info("send msg: ", msg)
}
