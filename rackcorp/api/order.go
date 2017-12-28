package api

import (
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

	var resp orderGetResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "OrderGet request failed for order Id '%s'.", orderId)
	}

	if resp.Code != "OK" || resp.Order == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Order, nil
}
