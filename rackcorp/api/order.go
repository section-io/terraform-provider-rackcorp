package api

import (
	"strconv"

	"github.com/pkg/errors"
)

type Storage struct {
	Name        string `json:"name,omitempty"`
	SizeMB      int    `json:"sizeMB"`
	StorageType string `json:"type"`
	SortOrder   int    `json:"order,omitempty"`
}

type ProductDetails struct {
	Credentials []Credential `json:"credentials"`
	Install     Install      `json:"install"`
	CpuCount    int          `json:"cpu"`
	Storage     []Storage    `json:"storage,omitempty"`
	MemoryGB    int          `json:"memoryGB"`
}

type Install struct {
	OperatingSystem   string `json:"operatingSystem"`
	PostInstallScript string `json:"postInstallScript,omitempty"`
	Template          string `json:"template,omitempty"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Order struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
	Status     string `json:"status"`
	ContractId string `json:"contractId"`
}

type ConfirmedOrder struct {
	ContractIds []string
}

type CreatedOrder struct {
	OrderId    string
	ChangeText string
}

type orderConfirmRequest struct {
	request
	OrderId string `json:"orderId"`
}

type orderConfirmResponse struct {
	response
	ContractIds []int `json:"contractID"`
}

type orderCreateRequest struct {
	request
	ProductCode    string         `json:"productCode"`
	CustomerId     string         `json:"customerId"`
	Quantity       int            `json:"quantity,omitempty"`
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

const (
	StorageTypeMagnetic = "MAGNETIC"
	StorageTypeSSD      = "SSD"
)

func sliceItoa(i []int) []string {
	a := make([]string, len(i))
	for index, value := range i {
		a[index] = strconv.Itoa(value)
	}
	return a
}

func GetVirtualServerProductCode(serverClass string, country string) string {
	return "SERVER_VIRTUAL_" + serverClass + "_" + country
}

func (c *client) OrderConfirm(orderId string) (*ConfirmedOrder, error) {
	if orderId == "" {
		return nil, errors.New("orderId parameter is required.")
	}

	req := &orderConfirmRequest{
		request: c.newRequest("order.confirm"),
		OrderId: orderId,
	}

	var resp orderConfirmResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "OrderConfirm request failed.")
	}

	if resp.Code != "OK" || len(resp.ContractIds) == 0 {
		return nil, newApiError(resp.response, nil)
	}

	return &ConfirmedOrder{
		ContractIds: sliceItoa(resp.ContractIds),
	}, nil
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
