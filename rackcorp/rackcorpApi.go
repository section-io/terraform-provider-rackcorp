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

	RackcorpApiDeviceGetCommand = "device.get"

	RackcorpApiOrderStatusPending  = "PENDING"
	RackcorpApiOrderStatusAccepted = "ACCEPTED"

	RackcorpApiOrderContractStatusActive  = "ACTIVE"
	RackcorpApiOrderContractStatusPending = "PENDING"
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
	DeviceId   string                   `json:"id"`
	Name       string                   `json:"name"`
	CustomerId string                   `json:"customerId"`
	PrimaryIP  string                   `json:"primaryIP"`
	Status     string                   `json:"status"`
	Extra      []RackcorpApiDeviceExtra `json:"extra"`
	// TODO assets, stdName, dateCreated, dateModified, dcDescription,
	//  dcId, cName, extra, firewallPolicies, ips, networkRoutes, ports,
	//  trafficCurrent, trafficEstimated, trafficMB, trafficShared
}

type RackcorpApiDeviceExtra struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type DeviceGetResponse struct {
	RackcorpApiResponse
	Device RackcorpApiDevice `json:"device"`
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
