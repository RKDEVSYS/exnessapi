package exnessapi

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint string
	Headers  *http.Header
	Message  chan WsTradeEvent
}

func newsWsConfig(endpoint string, headers *http.Header, message chan WsTradeEvent) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
		Message:  message,
		Headers:  headers,
	}
}

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	Dialer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}
	c, _, err := Dialer.Dial(cfg.Endpoint, *cfg.Headers)
	if err != nil {
		return nil, nil, err
	}
	c.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})

	/*
		This function will exit either on error from
		websocket.Conn.ReadMessage or then stopC channel
		is closed by the client
	*/
	go func() {
		defer close(doneC)
		if WebsocketKeepalive {
			keepAlive(c, WebsocketTimeout)
		}
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-cfg.Message:

			case <-doneC:
			}
			c.Close()
		}()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	return
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticket := time.NewTicker(timeout)
	lastResponse := time.Now()
	c.SetPongHandler(func(appData string) error {
		lastResponse = time.Now()
		return nil
	})
	go func() {
		deadline := time.Now().Add(10 * time.Second)
		err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
		if err != nil {
			return
		}
		<-ticket.C
		if time.Since(lastResponse) > timeout {
			c.Close()
			return
		}
	}()
}
