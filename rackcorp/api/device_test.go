package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestDeviceGet(t *testing.T) {
	defer gock.Off()

	const deviceId = "678"
	responseBody := getTestDataString(t, "device.get.responseBody.json")

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Post("/api/rest/v1/json.php").
		Reply(200).
		BodyString(responseBody)

	device, err := client.DeviceGet(deviceId)
	assert.Nil(t, err, "DeviceGet error")

	assert.Equal(t, "678", device.DeviceId, "DeviceId")

	assert.True(t, gock.IsDone(), "gock.IsDone")
}
