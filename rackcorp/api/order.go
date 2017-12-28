package api

import (
	"strconv"

	"github.com/pkg/errors"
)

type ProductDetails struct {
	Install  Install `json:"install"`
	CpuCount int     `json:"cpu"`
}

type Install struct {
	OperatingSystem string `json:"operatingSystem"`
}

type Order struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
	Status     string `json:"status"`
	ContractId string `json:"contractId"`
}

type CreatedOrder struct {
	OrderId    string
	ChangeText string
}

type orderCreateRequest struct {
	request
	ProductCode    string         `json:"productCode"`
	CustomerId     string         `json:"customerId"`
	ProductDetails ProductDetails `json:"productDetails"`
}

type orderCreateResponse struct {
	response
	OrderId    int    `json:"orderId"`
	ChangeText string `json:"changeTxt"`
	// TODO cost, currency, netCost, retailCost, retailNetCost
}

type orderGetRequest struct {
	request
	OrderId string `json:"orderId"`
}

type orderGetResponse struct {
	response
	Order *Order `json:"order"`
}

func (c *client) OrderCreate(productCode string, customerId string, productDetails ProductDetails) (*CreatedOrder, error) {
	if productCode == "" {
		return nil, errors.New("productCode parameter is required.")
	}

	if customerId == "" {
		return nil, errors.New("customerId parameter is required.")
	}

	req := &orderCreateRequest{
		request:        c.newRequest("order.create"),
		ProductCode:    productCode,
		CustomerId:     customerId,
		ProductDetails: productDetails,
	}

	var resp orderCreateResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "OrderCreate request failed.")
	}

	if resp.Code != "OK" || resp.OrderId == 0 {
		return nil, newApiError(resp.response, nil)
	}

	return &CreatedOrder{
		OrderId:    strconv.Itoa(resp.OrderId),
		ChangeText: resp.ChangeText,
	}, nil
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
