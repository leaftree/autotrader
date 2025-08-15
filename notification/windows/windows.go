package windows

import (
	"log"

	"github.com/gen2brain/beeep"
)

func SendNotify(msg string) {
	go func() {
		beeep.AppName = "autotrader"
		err := beeep.Notify("GateIO 交易通知", "ETH_USDT 多单下单成功\n价格：1\nstop loss: 2 ", "xu.jpeg")
		if err != nil {
			log.Println("send notify failed: ", err)
			return
		}

		/*
			streamer, format, err := wav.Decode(strings.NewReader("msg.wav"))
			if err != nil {
				log.Panicln(err)
				return
			}
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			defer speaker.Close()
			done := make(chan bool)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() { done <- true })))
			<-done // 等待播放完成
		*/
	}()
}
