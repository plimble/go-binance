package binance

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	baseURL = "wss://stream2.binance.com:9443/ws"
)

// WsDepthHandler handle websocket depth event
type WsDiffDepthHandler func(event *WsDiffDepthEvent)
type WsPartialBookDepthHandler func(event *WsPartialBookDepthEvent)

// WsPartialBookDepthServe Top <levels> bids and asks, pushed every second. Valid <levels> are 5, 10, or 20.
func WsPartialBookDepthServe(symbol string, levels string, handler WsPartialBookDepthHandler, errHandler WsErrorHandler) *WsService {
	endpoint := fmt.Sprintf("%s/%s@depth%s", baseURL, strings.ToLower(symbol), levels)
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			// TODO: callback if there is an error
			return
		}
		event := new(WsPartialBookDepthEvent)
		event.LastUpdateId = j.Get("lastUpdateId").MustInt64()
		bidsLen := len(j.Get("bids").MustArray())
		event.Bids = make([]Bid, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := j.Get("bids").GetIndex(i)
			event.Bids[i] = Bid{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		asksLen := len(j.Get("asks").MustArray())
		event.Asks = make([]Ask, asksLen)
		for i := 0; i < asksLen; i++ {
			item := j.Get("asks").GetIndex(i)
			event.Asks[i] = Ask{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		handler(event)
	}

	return newWsService(endpoint, wsHandler, errHandler)
}

// WsPartialBookDepthEvent define websocket partial orderbook depth event
type WsPartialBookDepthEvent struct {
	LastUpdateId int64 `json:"lastUpdateId"`
	Bids         []Bid `json:"bids"`
	Asks         []Ask `json:"asks"`
}

// WsDiffDepthServe Order book price and quantity depth updates used to locally manage an order book pushed every second.
func WsDiffDepthServe(symbol string, handler WsDiffDepthHandler, errHandler WsErrorHandler) *WsService {
	endpoint := fmt.Sprintf("%s/%s@depth", baseURL, strings.ToLower(symbol))
	wsHandler := func(message []byte) {
		j, err := newJSON(message)
		if err != nil {
			// TODO: callback if there is an error
			return
		}
		event := new(WsDiffDepthEvent)
		event.Event = j.Get("e").MustString()
		event.Time = j.Get("E").MustInt64()
		event.Symbol = j.Get("s").MustString()
		event.UpdateID = j.Get("u").MustInt64()
		bidsLen := len(j.Get("b").MustArray())
		event.Bids = make([]Bid, bidsLen)
		for i := 0; i < bidsLen; i++ {
			item := j.Get("b").GetIndex(i)
			event.Bids[i] = Bid{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		asksLen := len(j.Get("a").MustArray())
		event.Asks = make([]Ask, asksLen)
		for i := 0; i < asksLen; i++ {
			item := j.Get("a").GetIndex(i)
			event.Asks[i] = Ask{
				Price:    item.GetIndex(0).MustString(),
				Quantity: item.GetIndex(1).MustString(),
			}
		}
		handler(event)
	}

	return newWsService(endpoint, wsHandler, errHandler)
}

// WsDepthEvent define websocket depth event
type WsDiffDepthEvent struct {
	Event    string `json:"e"`
	Time     int64  `json:"E"`
	Symbol   string `json:"s"`
	UpdateID int64  `json:"u"`
	Bids     []Bid  `json:"b"`
	Asks     []Ask  `json:"a"`
}

// WsKlineHandler handle websocket kline event
type WsKlineHandler func(event *WsKlineEvent)

// WsKlineServe serve websocket kline handler with a symbol and interval like 15m, 30s
func WsKlineServe(symbol string, interval string, handler WsKlineHandler, errHandler WsErrorHandler) *WsService {
	endpoint := fmt.Sprintf("%s/%s@kline_%s", baseURL, strings.ToLower(symbol), interval)
	wsHandler := func(message []byte) {
		event := new(WsKlineEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			return
		}
		handler(event)
	}
	return newWsService(endpoint, wsHandler, errHandler)
}

// WsKlineEvent define websocket kline event
type WsKlineEvent struct {
	Event  string  `json:"e"`
	Time   int64   `json:"E"`
	Symbol string  `json:"s"`
	Kline  WsKline `json:"k"`
}

// WsKline define websocket kline
type WsKline struct {
	StartTime            int64  `json:"t"`
	EndTime              int64  `json:"T"`
	Symbol               string `json:"s"`
	Interval             string `json:"i"`
	FirstTradeID         int64  `json:"f"`
	LastTradeID          int64  `json:"L"`
	Open                 string `json:"o"`
	Close                string `json:"c"`
	High                 string `json:"h"`
	Low                  string `json:"l"`
	Volume               string `json:"v"`
	TradeNum             int64  `json:"n"`
	IsFinal              bool   `json:"x"`
	QuoteVolume          string `json:"q"`
	ActiveBuyVolume      string `json:"V"`
	ActiveBuyQuoteVolume string `json:"Q"`
}

// WsAggTradeHandler handle websocket aggregate trade event
type WsAggTradeHandler func(event *WsAggTradeEvent)

// WsAggTradeServe serve websocket aggregate handler with a symbol
func WsAggTradeServe(symbol string, handler WsAggTradeHandler, errHandler WsErrorHandler) *WsService {
	endpoint := fmt.Sprintf("%s/%s@aggTrade", baseURL, strings.ToLower(symbol))
	wsHandler := func(message []byte) {
		event := new(WsAggTradeEvent)
		err := json.Unmarshal(message, event)
		if err != nil {
			return
		}
		handler(event)
	}

	return newWsService(endpoint, wsHandler, errHandler)
}

// WsAggTradeEvent define websocket aggregate trade event
type WsAggTradeEvent struct {
	Event                 string `json:"e"`
	Time                  int64  `json:"E"`
	Symbol                string `json:"s"`
	AggTradeID            int64  `json:"a"`
	Price                 string `json:"p"`
	Quantity              string `json:"q"`
	FirstBreakdownTradeID int64  `json:"f"`
	LastBreakdownTradeID  int64  `json:"l"`
	TradeTime             int64  `json:"T"`
	IsBuyerMaker          bool   `json:"m"`
	Placeholder           bool   `json:"M"` // add this field to avoid case insensitive unmarshaling
}

// WsUserDataServe serve user data handler with listen key
func WsUserDataServe(listenKey string, handler WsHandler, errHandler WsErrorHandler) *WsService {
	endpoint := fmt.Sprintf("%s/%s", baseURL, listenKey)
	return newWsService(endpoint, handler, errHandler)
}
