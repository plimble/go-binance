package binance

import (
	"github.com/gorilla/websocket"
)

type WsService struct {
	endpoint   string
	close      chan struct{}
	handler    WsHandler
	errHandler WsErrorHandler
}

func newWsService(endpoint string, handler WsHandler, errHandler WsErrorHandler) *WsService {
	if handler == nil {
		handler = defaultWsHandler
	}

	if errHandler == nil {
		errHandler = defaultWsErrorHandler
	}

	return &WsService{
		endpoint:   endpoint,
		close:      make(chan struct{}, 1),
		handler:    handler,
		errHandler: errHandler,
	}
}

func (w *WsService) Close() {
	w.close <- struct{}{}
}

func (w *WsService) Serve() error {
	c, _, err := websocket.DefaultDialer.Dial(w.endpoint, nil)
	if err != nil {
		return err
	}

	defer c.Close()
	for {
		select {
		case <-w.close:
			c.Close()
			return nil
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				w.errHandler(err)
			} else {
				w.handler(message)
			}
		}
	}
}

// WsHandler handle raw websocket message
type WsHandler func(message []byte)
type WsErrorHandler func(err error)

var defaultWsErrorHandler = func(err error) {
	panic(err)
}

var defaultWsHandler = func(message []byte) {}
