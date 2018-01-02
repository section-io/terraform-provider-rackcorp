package api

import (
	"strconv"

	"github.com/pkg/errors"
)

type createdTransaction struct {
	TransactionId        int    `json:"rcTransactionId"`
	ConfirmationRequired bool   `json:"confirmationRequired"`
	ConfirmationText     string `json:"confirmationText"`
	ObjectType           string `json:"objType"`
	ObjectId             string `json:"objId"`
	Type                 string `json:"type"`
	Data                 string `json:"data"`
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
	Transaction *createdTransaction `json:"rcTransaction"`
}

type existingTransaction struct {
	TransactionId string `json:"rcTransactionId"`
	ObjectType    string `json:"objType"`
	ObjectId      string `json:"objId"`
	Type          string `json:"method"`
	Data          string `json:"data"`
	Status        string `json:"status"`
	StatusInfo    string `json:"statusInfo"`
}

type transactionGetRequest struct {
	request
	TransactionId int `json:"rcTransactionId"`
}

type transactionGetResponse struct {
	response
	Transaction *existingTransaction `json:"rcTransaction"`
}

type Transaction struct {
	TransactionId        string
	ObjectType           string
	ObjectId             string
	Type                 string
	Data                 string
	ConfirmationRequired bool
	ConfirmationText     string
	Status               string
	StatusInfo           string
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

func (t *createdTransaction) ToTransaction() *Transaction {
	return &Transaction{
		TransactionId:        strconv.Itoa(t.TransactionId),
		ObjectType:           t.ObjectType,
		ObjectId:             t.ObjectId,
		Type:                 t.Type,
		Data:                 t.Data,
		ConfirmationRequired: t.ConfirmationRequired,
		ConfirmationText:     t.ConfirmationText,
	}
}

func (t *existingTransaction) ToTransaction() *Transaction {
	return &Transaction{
		TransactionId: t.TransactionId,
		ObjectType:    t.ObjectType,
		ObjectId:      t.ObjectId,
		Type:          t.Type,
		Data:          t.Data,
		Status:        t.Status,
		StatusInfo:    t.StatusInfo,
	}
}

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

	return resp.Transaction.ToTransaction(), nil
}

func (c *client) TransactionGet(transactionId string) (*Transaction, error) {
	numericId, err := strconv.Atoi(transactionId)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid transactionId")
	}

	req := &transactionGetRequest{
		request:       c.newRequest("rctransaction.get"),
		TransactionId: numericId,
	}

	var resp transactionGetResponse
	err = c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "TransactionCreate request failed.")
	}

	if resp.Code != "OK" || resp.Transaction == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Transaction.ToTransaction(), nil
}
