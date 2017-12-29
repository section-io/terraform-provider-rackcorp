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
}

type transactionCreateResponse struct {
	response
	Transaction *Transaction `json:"rcTransaction"`
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

func (c *client) TransactionCreate(transactionType string, objectType string, objectId string) (*Transaction, error) {
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
