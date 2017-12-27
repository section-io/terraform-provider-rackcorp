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

	RackcorpApiOrderCreateCommand  = "order.create"
	RackcorpApiOrderConfirmCommand = "order.confirm"
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
	ContractId []int `json:"contractID"`
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
