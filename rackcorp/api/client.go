package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type request struct {
	ApiUuid   string `json:"APIUUID"`
	ApiSecret string `json:"APISECRET"`
	Command   string `json:"cmd"`
}

type response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Debug   string `json:"debug"`
}

type client struct {
	address string
	uuid    string
	secret  string
}

type Client interface {
	OrderConfirm(orderId string) (*ConfirmedOrder, error)
	OrderCreate(productCode string, customerId string, productDetails ProductDetails) (*CreatedOrder, error)
	OrderGet(orderId string) (*Order, error)

	OrderContractGet(contractId string) (*OrderContract, error)

	DeviceGet(deviceId string) (*Device, error)

	TransactionCreate(transactionType string, objectType string, objectId string, confirm bool) (*Transaction, error)
	TransactionGet(transactionId string) (*Transaction, error)
	TransactionGetAll(filter TransactionFilter) ([]Transaction, int, error)
}

const (
	DefaultAddress = "https://api.rackcorp.net/api/rest/v1/json.php"
)

func NewClient(uuid string, secret string) (Client, error) {
	if uuid == "" {
		return nil, errors.New("uuid argument must not be empty.")
	}

	if secret == "" {
		return nil, errors.New("secret argument must not be empty.")
	}

	return &client{
		address: DefaultAddress,
		uuid:    uuid,
		secret:  secret,
	}, nil
}

func safeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func (c client) newRequest(command string) request {
	if command == "" {
		panic("command is required.")
	}

	return request{
		ApiUuid:   c.uuid,
		ApiSecret: c.secret,
		Command:   command,
	}
}

func (c client) httpPost(requestBody []byte) (responseBody []byte, outErr error) {
	response, err := http.Post(c.address, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, errors.Wrap(err, "HTTP POST failed for request.")
	}
	defer safeClose(response.Body, &outErr)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c client) httpPostJson(req interface{}, resp interface{}) error {

	reqBody, err := json.Marshal(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to JSON encode request: %v", req)
	}

	respBody, err := c.httpPost(reqBody)
	if err != nil {
		return errors.Wrapf(err, "Failed to HTTP POST request: %v", req)
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return errors.Wrapf(err, "Could not JSON decode response: %s", respBody)
	}

	return nil
}
