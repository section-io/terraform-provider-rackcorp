package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Order struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
	Status     string `json:"status"`
	ContractId string `json:"contractId"`
}

type orderGetRequest struct {
	request
	OrderId string `json:"orderId"`
}

type orderGetResponse struct {
	response
	Order *Order `json:"order"`
}

func (c *client) OrderGet(orderId string) (*Order, error) {
	if orderId == "" {
		return nil, errors.New("orderId parameter is required.")
	}

	req := &orderGetRequest{
		request: c.newRequest("order.get"),
		OrderId: orderId,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", req)
	}

	respBody, err := c.httpPost(reqBody) // TODO c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", req)
	}

	var resp orderGetResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", respBody)
	}

	if resp.Code != "OK" || resp.Order == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Order, nil
}
