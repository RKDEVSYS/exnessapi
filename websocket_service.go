package exnessapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RKDEVSYS/exnessapi/types"
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

type WsSubscriberHandler func(event *types.WsTradeEvent)

func WsSubscriber(symbol string, handler WsSubscriberHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	subscriberMessage := make(chan types.WsTradeEvent, 1)
	subscriberMessage <- types.WsTradeEvent{Type: "TicksSubscribe", Body: &types.wsTradeBody{Symbol: symbol}}
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
