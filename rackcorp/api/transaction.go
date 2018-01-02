package api

import (
	"github.com/pkg/errors"
)

type Transaction struct {
	TransactionId        int    `json:"rcTransactionId"`
	ConfirmationRequired bool   `json:"confirmationRequired"`
	ConfirmationText     string `json:"confirmationText"`
	ObjectType           string `json:"objType"`
	ObjectId             string `json:"objId"`
	Type                 string `json:"type"`
	// TODO data when type is known
}

type transactionCreateRequest struct {
	request
	ObjectType string `json:"objType"`
	ObjectId   string `json:"objId"`
	Type       string `json:"type"`
	Confirm    bool   `json:"confirmation"`
}

type transactionCreateResponse struct {
	response
	Transaction *Transaction `json:"rcTransaction"`
}

type TransactionGet struct {
	Data          string `json:"data"`
	Method        string `json:"method"`
	ObjectId      string `json:"objId"`
	ObjectType    string `json:"objType"`
	TransactionId string `json:"rcTransactionId"`
	Status        string `json:"status"`
	StatusInfo    string `json:"statusInfo"`
}

type transactionGetRequest struct {
	request
	TransactionId int `json:"rcTransactionId"`
}

type transactionGetResponse struct {
	response
	Transaction *TransactionGet `json:"rcTransaction"`
}

const (
	TransactionObjectTypeDevice = "DEVICE"

	TransactionTypeCancel        = "CANCEL"
	TransactionTypeCloseVNC      = "CLOSEVNC"
	TransactionTypeForceShutdown = "FORCESHUTDOWN"
	TransactionTypeOpenVNC       = "OPENVNC" // data parameter contains public IP that allows VNC
	TransactionTypeRefreshConfig = "REFRESHCONFIG"
	TransactionTypeSafeShutdown  = "SAFESHUTDOWN"
	TransactionTypeShutdown      = "SHUTDOWN"
	TransactionTypeStartup       = "STARTUP"
)

func (c *client) TransactionCreate(transactionType string, objectType string, objectId string, confirm bool) (*Transaction, error) {
	if transactionType == "" {
		return nil, errors.New("transactionType parameter is required.")
	}

	if objectType == "" {
		return nil, errors.New("objectType parameter is required.")
	}

	if objectId == "" {
		return nil, errors.New("objectId parameter is required.")
	}

	req := &transactionCreateRequest{
		request:    c.newRequest("rctransaction.create"),
		Type:       transactionType,
		ObjectType: objectType,
		ObjectId:   objectId,
		Confirm:    confirm,
	}

	var resp transactionCreateResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "TransactionCreate request failed.")
	}

	if resp.Code != "OK" || resp.Transaction == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Transaction, nil
}

func (c *client) TransactionGet(transactionId int) (*TransactionGet, error) {
	req := &transactionGetRequest{
		request:       c.newRequest("rctransaction.get"),
		TransactionId: transactionId,
	}

	var resp transactionGetResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "TransactionCreate request failed.")
	}

	if resp.Code != "OK" || resp.Transaction == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Transaction, nil
}
