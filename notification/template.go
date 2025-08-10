package notification

import (
	"fmt"
	"text/template"
)

var (
	CreateOrderTemplateText = "create {{.side}} order, price: {{.price}}, size: {{.size}}, stop loss price: {{.stop_loss_price}}"
)

var t *template.Template

func init() {
	var err error
	t, err = template.New("trader").Parse(CreateOrderTemplateText)
	if err != nil {
		fmt.Println(err)
	}
}
