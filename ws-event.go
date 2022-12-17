package exnessapi

type WsTradeEvent struct {
	Type    string       `json:"type,omitempty"`
	Body    *WsTradeBody `json:"body,omitempty"`
	Channel uint         `json:"channel,omitempty"`
	Id      uint         `json:"id,omitempty"`
}

type WsTradeBody struct {
	Type      interface{} `json:"type,omitempty"`
	Symbol    string      `json:"symbol,omitempty"`
	Price     interface{} `json:"price,omitempty"`
	Volume    interface{} `json:"volume,omitempty"`
	Deviation interface{} `json:"deviation,omitempty"`
	Sl        interface{} `json:"sl,omitempty"`
	Tp        interface{} `json:"tp,omitempty"`
	Time      interface{} `json:"time,omitempty"`
	Bid       interface{} `json:"bid,omitempty"`
	Ask       interface{} `json:"ask,omitempty"`
}
