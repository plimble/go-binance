package binance

import (
	"github.com/gorilla/websocket"
)

type WsService struct {
	endpoint   string
	close      chan struct{}
	handler    WsHandler
	errHandler WsErrorHandler
	c          *websocket.Conn
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

func (w *WsService) Connect() error {
	var err error
	w.c, _, err = websocket.DefaultDialer.Dial(w.endpoint, nil)

	return err
}

func (w *WsService) Serve() {
	defer w.c.Close()
	for {
		select {
		case <-w.close:
			w.c.Close()
			return
		default:
			_, message, err := w.c.ReadMessage()
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
