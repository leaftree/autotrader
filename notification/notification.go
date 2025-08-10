package notification

import (
	"bytes"
	"context"
	"fmt"
	"log"

	wxpusher "github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

const (
	appToken = "AT_4EhDwiLfZfre2eYGWRfiPoeFkNlbciIW"
	userID   = "UID_XK7Qp5fAPTxNgszplAEAqgnOTebX"
)

func SendNotify(msg string) {
	content := model.NewMessage(appToken).SetContent(msg).AddUId(userID)
	resp, err := wxpusher.SendMessage(content)
	if err != nil {
		fmt.Println(resp, err)
	}
}

func SendCreateOrderNotification(ctx context.Context, values map[string]any) {
	var b bytes.Buffer
	err := t.Execute(b, values)
	if err != nil {
		log.Println("generate create order notification failed: ", err)
		return
	}
	SendNotify(b.String())
}
