package binance

import (
	"context"
	"fmt"
	"testing"
)

func TestExchangeInfoService(t *testing.T) {
	c := NewClient("", "")
	req := c.NewExchangeInfoService()
	res, err := req.Do(context.Background())

	if err != nil {
		t.Error(err)
	}

	fmt.Println(res)
}
