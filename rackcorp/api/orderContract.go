package api

import (
	"github.com/pkg/errors"
)

type OrderContract struct {
	ContractId string `json:"contractId"`
	CustomerId string `json:"customerId"`
	DeviceId   string `json:"deviceID"`
	Status     string `json:"status"` // TODO enum
	Type       string `json:"type"`   // TODO enum
	// TODO contractInfo, created, currency, lastBilled, modified, notes, referenceID, serviceBillId
}

type orderContractGetRequest struct {
	request
	ContractId string `json:"contractId"`
}

type orderContractGetResponse struct {
	response
	Contract *OrderContract `json:"contract"`
}

func (c *client) OrderContractGet(contractId string) (*OrderContract, error) {
	if contractId == "" {
		return nil, errors.New("contractId parameter is required.")
	}

	req := &orderContractGetRequest{
		request:    c.newRequest("order.contract.get"),
		ContractId: contractId,
	}

	var resp orderContractGetResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "OrderContractGet request failed for contract Id '%s'.", contractId)
	}

	if resp.Code != "OK" || resp.Contract == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Contract, nil
}
