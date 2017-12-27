package rackcorp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	RackcorpApiResponseCodeOK           = "OK"
	RackcorpApiResponseCodeAccessDenied = "ACCESS_DENIED"
	RackcorpApiResponseCodeFault        = "FAULT"

	RackcorpApiDeviceGetCommand        = "device.get"
	RackcorpApiOrderConfirmCommand     = "order.confirm"
	RackcorpApiOrderContractGetCommand = "order.contract.get"
	RackcorpApiOrderCreateCommand      = "order.create"
	RackcorpApiOrderGetCommand         = "order.get"

	RackcorpApiOrderStatusPending  = "PENDING"
	RackcorpApiOrderStatusAccepted = "ACCEPTED"

	RackcorpApiOrderContractStatusActive  = "ACTIVE"
	RackcorpApiOrderContractStatusPending = "PENDING"

	RackcorpApiOrderContractTypeVirtualServer = "VIRTUALSERVER"
)

type RackcorpApiRequest struct {
	ApiUuid   string `json:"APIUUID"`
	ApiSecret string `json:"APISECRET"`
	Command   string `json:"cmd"`
}

type RackcorpApiResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OrderGetRequest struct {
	RackcorpApiRequest
	OrderId string `json:"orderId"`
}

func NewOrderGetRequest(orderId string) *OrderGetRequest {
	return &OrderGetRequest{
		RackcorpApiRequest: RackcorpApiRequest{
			Command: RackcorpApiOrderGetCommand,
		},
		OrderId: orderId,
	}
}

type RackcorpApiOrder struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
	Status     string `json:"status"`
	ContractId string `json:"contractId"`
}

type OrderGetResponse struct {
	RackcorpApiResponse
	Order RackcorpApiOrder `json:"order"`
}

type OrderContractGetRequest struct {
	RackcorpApiRequest
	ContractId string `json:"contractId"`
}

func NewOrderContractGetRequest(contractId string) *OrderContractGetRequest {
	return &OrderContractGetRequest{
		RackcorpApiRequest: RackcorpApiRequest{
			Command: RackcorpApiOrderContractGetCommand,
		},
		ContractId: contractId,
	}
}

type RackcorpApiContract struct {
	ContractId string `json:"contractId"`
	CustomerId string `json:"customerId"`
	DeviceId   string `json:"deviceID"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	// TODO contractInfo, created, currency, lastBilled, modified, notes, referenceID, serviceBillId

}

type OrderContractGetResponse struct {
	RackcorpApiResponse
	Contract RackcorpApiContract `json:"contract"`
}

type DeviceGetRequest struct {
	RackcorpApiRequest
	DeviceId string `json:"deviceId"`
}

func NewDeviceGetRequest(deviceId string) *DeviceGetRequest {
	return &DeviceGetRequest{
		RackcorpApiRequest: RackcorpApiRequest{
			Command: RackcorpApiDeviceGetCommand,
		},
		DeviceId: deviceId,
	}
}

type RackcorpApiDevice struct {
	DeviceId   string `json:"id"`
	Name       string `json:"name"`
	CustomerId string `json:"customerId"`
	PrimaryIP  string `json:"primaryIP"`
	Status     string `json:"status"`
	// TODO assets, stdName, dateCreated, dateModified, dcDescription,
	//  dcId, cName, extra, firewallPolicies, ips, networkRoutes, ports,
	//  trafficCurrent, trafficEstimated, trafficMB, trafficShared
}

type DeviceGetResponse struct {
	RackcorpApiResponse
	Device RackcorpApiDevice `json:"device"`
}

type OrderCreateRequest struct {
	RackcorpApiRequest
	ProductCode    string         `json:"productCode"`
	CustomerId     string         `json:"customerId"`
	ProductDetails ProductDetails `json:"productDetails"`
}

func NewOrderCreateRequest() *OrderCreateRequest {
	return &OrderCreateRequest{
		RackcorpApiRequest: RackcorpApiRequest{
			Command: RackcorpApiOrderCreateCommand,
		},
	}
}

type ProductDetails struct {
	Install Install `json:"install"`
}

type Install struct {
	OperatingSystem string `json:"operatingSystem"`
}

type OrderCreateResponse struct {
	RackcorpApiResponse
	OrderId    int    `json:"orderId"`
	ChangeText string `json:"changeTxt"`
	// TODO cost, currency, netCost, retailCost, retailNetCost
}

type OrderConfirmRequest struct {
	RackcorpApiRequest
	OrderId string `json:"orderId"`
}

func NewOrderConfirmRequest(orderId string) *OrderConfirmRequest {
	return &OrderConfirmRequest{
		RackcorpApiRequest: RackcorpApiRequest{
			Command: RackcorpApiOrderConfirmCommand,
		},
		OrderId: orderId,
	}
}

type OrderConfirmResponse struct {
	RackcorpApiResponse
	ContractIds []int `json:"contractID"`
}

func safeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func (request *RackcorpApiRequest) Configure(config Config) {
	request.ApiUuid = config.ApiUuid
	request.ApiSecret = config.ApiSecret // TODO exclude from logs
}

func postRackcorpApiRequest(requestBody []byte, config Config) (responseBody []byte, outErr error) {

	response, err := http.Post(config.ApiAddress, "application/json", bytes.NewReader(requestBody))
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

func (request *OrderCreateRequest) Post(config Config) (*OrderCreateResponse, error) {
	request.RackcorpApiRequest.Configure(config)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", request)
	}

	responseBody, err := postRackcorpApiRequest(requestBody, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", request)
	}

	var response OrderCreateResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", responseBody)
	}

	return &response, nil
}

func (request *OrderConfirmRequest) Post(config Config) (*OrderConfirmResponse, error) {
	request.RackcorpApiRequest.Configure(config)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", request)
	}

	responseBody, err := postRackcorpApiRequest(requestBody, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", request)
	}

	var response OrderConfirmResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", responseBody)
	}

	return &response, nil
}

func (request *OrderGetRequest) Post(config Config) (*OrderGetResponse, error) {
	request.RackcorpApiRequest.Configure(config)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", request)
	}

	responseBody, err := postRackcorpApiRequest(requestBody, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", request)
	}

	var response OrderGetResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", responseBody)
	}

	return &response, nil
}

func (request *OrderContractGetRequest) Post(config Config) (*OrderContractGetResponse, error) {
	request.RackcorpApiRequest.Configure(config)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", request)
	}

	responseBody, err := postRackcorpApiRequest(requestBody, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", request)
	}

	var response OrderContractGetResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", responseBody)
	}

	return &response, nil
}

func (request *DeviceGetRequest) Post(config Config) (*DeviceGetResponse, error) {
	request.RackcorpApiRequest.Configure(config)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON encode request: %v", request)
	}

	responseBody, err := postRackcorpApiRequest(requestBody, config)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to HTTP POST request: %v", request)
	}

	var response DeviceGetResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not JSON decode response: %s", responseBody)
	}

	return &response, nil
}
