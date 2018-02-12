package binance

import (
	"context"
	"encoding/json"
)

// ExchangeInfoService show exchange info
type ExchangeInfoService struct {
	c *Client
}

// Do send request
func (s *ExchangeInfoService) Do(ctx context.Context, opts ...RequestOption) (*ExchangeInfoResponse, error) {
	r := &request{
		method:   "GET",
		endpoint: "/api/v1/exchangeInfo",
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := new(ExchangeInfoResponse)

	if err = json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	return res, nil
}

// DepthResponse define depth info with bids and asks
type ExchangeInfoResponse struct {
	Timezone   string                   `json:"timezone"`
	ServerTime int64                    `json:"serverTime"`
	RateLimits []*ExchangeInfoRateLimit `json:"rateLimits"`
	Symbols    []*ExchangeInfoSymbol    `json:"symbols"`
}

type ExchangeInfoRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	Limit         int    `json:"limit"`
}

type ExchangeInfoSymbol struct {
	Symbol             string                `json:"symbol"`
	Status             string                `json:"status"`
	BaseAsset          string                `json:"baseAsset"`
	BaseAssetPrecision int                   `json:"baseAssetPrecision"`
	QuoteAsset         string                `json:"quoteAsset"`
	QuotePrecision     int                   `json:"quotePrecision"`
	OrderTypes         []string              `json:"orderTypes"`
	IcebergAllowed     bool                  `json:"icebergAllowed"`
	Filter             []*ExchangeInfoFilter `json:"filters"`
}

type ExchangeInfoFilter struct {
	FilterType  string `json:"filterType"`
	MinPrice    string `json:"minPrice"`
	MaxPrice    string `json:"maxPrice"`
	TickSize    string `json:"tickSize"`
	MinQty      string `json:"minQty"`
	MaxQty      string `json:"maxQty"`
	StepSize    string `json:"stepSize"`
	MinNotional string `json:"minNotional"`
}
