package binance

import (
	"context"
	"encoding/json"
)

type SystemStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

// SystemStatusService check binance sever status
type SystemStatusService struct {
	c *Client
}

// Do send request
func (s *SystemStatusService) Do(ctx context.Context) (*SystemStatus, error) {
	r := &request{
		method:   "GET",
		endpoint: "/wapi/v3/systemStatus.html",
		// secType:  secTypeSigned,
	}

	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return nil, err
	}

	res := &SystemStatus{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
