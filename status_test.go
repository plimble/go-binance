package binance

import (
	"context"
	"fmt"
	"testing"
)

func TestSystemStatusService(t *testing.T) {
	c := NewClient("", "")
	s := c.NewSystemStatusService()
	res, err := s.Do(context.Background())
	fmt.Println(err)
	fmt.Println(res)
}
