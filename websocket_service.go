package exnessapi

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	baseWsMainUrl = "wss://my.exness.com/wts-api"
	wsAccountPath = "trial11"
)

var (
	WebsocketTimeout   = 60 * time.Second
	WebsocketKeepalive = false
)

func getWsEndpoint() string {
	return fmt.Sprintf("%s/%s", baseWsMainUrl, wsAccountPath)
}

type WsSubscriberHandler func(event *WsTradeEvent)

func WsSubscriber(symbol string, handler WsSubscriberHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	subscriberMessage := make(chan WsTradeEvent, 1)
	subscriberMessage <- WsTradeEvent{Type: "TicksSubscribe", Body: &wsTradeBody{Symbol: symbol}}
	cfg := newsWsConfig(getWsEndpoint(), subscriberMessage)
	wsHandler := func(message []byte) {
		event := new(WsTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(event)
	}
	return wsServe(cfg, wsHandler, errHandler)
}
