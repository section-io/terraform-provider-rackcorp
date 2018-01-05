package api

import (
	"github.com/pkg/errors"
)

type DeviceExtra struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Device struct {
	DeviceId     string        `json:"id"`
	Name         string        `json:"name"`
	CustomerId   string        `json:"customerId"`
	PrimaryIP    string        `json:"primaryIP"`
	Status       string        `json:"status"`
	Extra        []DeviceExtra `json:"extra"`
	DataCenterId string        `json:"dcid"`
	// TODO assets, stdName, dateCreated, dateModified, dcDescription,
	//  dcName, firewallPolicies, ips, networkRoutes, ports,
	//  trafficCurrent, trafficEstimated, trafficMB, trafficShared
}

type deviceGetRequest struct {
	request
	DeviceId string `json:"deviceId"`
}

type deviceGetResponse struct {
	response
	Device *Device `json:"device"`
}

func (c *client) DeviceGet(deviceId string) (*Device, error) {
	if deviceId == "" {
		return nil, errors.New("deviceId parameter is required.")
	}

	req := &deviceGetRequest{
		request:  c.newRequest("device.get"),
		DeviceId: deviceId,
	}

	var resp deviceGetResponse
	err := c.httpPostJson(req, &resp)
	if err != nil {
		return nil, errors.Wrapf(err, "DeviceGet request failed for device Id '%s'.", deviceId)
	}

	if resp.Code != "OK" || resp.Device == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Device, nil
}
