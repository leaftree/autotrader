package notification

import (
	"bytes"
	"context"
	"runtime"

	"github.com/leaftree/autotrader/config"
	log "github.com/leaftree/autotrader/logger"
	"github.com/leaftree/autotrader/notification/windows"
	"github.com/leaftree/autotrader/notification/wxpusher"
)

var (
	logger = log.NewLogger("notification")
)

func SendCreateOrderNotification(ctx context.Context, values map[string]any) {
	var b bytes.Buffer
	err := t.Execute(&b, values)
	if err != nil {
		logger.Errorf("generate create order notification failed: ", err)
		return
	}
	wxpusher.SendNotify(b.String())

	if runtime.GOOS == "windows" && config.GetConfig().Notification.Windows == "on" {
		windows.SendNotify(b.String())
	}
}
