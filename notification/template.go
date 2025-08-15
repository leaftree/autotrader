package notification

import (
	"fmt"
	"text/template"
)

var (
	//<p><strong style="color: red;">这是一段加粗且红色的文字</strong></p>
	//CreateOrderTemplateText = "create ***{{.contract}}*** {{.side}} order, \nprice: {{.price}}, size: {{.size}}, stop loss price: {{.stop_loss_price}}"
	CreateOrderTemplateText = "create <strong style=\"color:red;\">{{.contract}}</strong> {{.side}} order\nprice: {{.price}}\nsize: {{.size}}\nstop loss price: {{.stop_loss_price}}"
)

var t *template.Template

func init() {
	var err error
	t, err = template.New("trader").Parse(CreateOrderTemplateText)
	if err != nil {
		fmt.Println(err)
	}
}
