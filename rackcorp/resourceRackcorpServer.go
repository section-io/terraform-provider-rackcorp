package rackcorp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func safeClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func resourceRackcorpServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceRackcorpServerCreate,
		Delete: resourceRackcorpServerDelete,
		Read:   resourceRackcorpServerRead,
		Schema: map[string]*schema.Schema{
			"country": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_class": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) (outErr error) {
	config := meta.(Config)

	install := Install{
		OperatingSystem: d.Get("server_class").(string),
	}

	productDetails := ProductDetails{
		Install: install,
	}

	orderRequest := OrderRequest{
		ApiUuid:        config.ApiUuid,
		ApiSecret:      config.ApiSecret,
		Command:        "order.create",
		CustomerId:     config.CustomerId,
		ProductCode:    "SERVER_VIRTUAL_" + d.Get("server_class").(string) + "_" + d.Get("country").(string),
		ProductDetails: productDetails,
	}

	orderRequestJson, err := json.Marshal(orderRequest)
	if err != nil {
		panic(err)
	}

	order, err := http.Post(config.ApiAddress, "application/json", bytes.NewBuffer(orderRequestJson))
	if err != nil {
		panic(err)
	}

	var orderResponse OrderResponse
	decodeErr := json.NewDecoder(order.Body).Decode(&orderResponse)
	if decodeErr != nil {
		panic(decodeErr)
	}

	defer safeClose(order.Body, &outErr)

	if orderResponse.Code != "OK" {
		panic(orderResponse.Code)
	}

	confirmRequest := ConfirmRequest{
		ApiUuid:    config.ApiUuid,
		ApiSecret:  config.ApiSecret,
		Command:    "order.confirm",
		CustomerId: config.CustomerId,
		OrderId:    strconv.Itoa(orderResponse.OrderId),
	}

	confirmRequestJson, err := json.Marshal(confirmRequest)
	if err != nil {
		panic(err)
	}

	confirm, err := http.Post(config.ApiAddress, "application/json", bytes.NewBuffer(confirmRequestJson))
	var confirmResponse ConfirmResponse
	decodeErr = json.NewDecoder(confirm.Body).Decode(&confirmResponse)
	if decodeErr != nil {
		panic(decodeErr)
	}

	if confirmResponse.Code != "OK" {
		panic(fmt.Sprintf("%#v\n", confirmResponse))
	}

	// panic(fmt.Sprintf("%#v\n", confirmResponse))

	defer safeClose(confirm.Body, &outErr)

	return resourceRackcorpServerRead(d, meta)
}

func resourceRackcorpServerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackcorpServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

type OrderRequest struct {
	ApiUuid        string         `json:"APIUUID"`
	ApiSecret      string         `json:"APISECRET"`
	Command        string         `json:"cmd"`
	ProductCode    string         `json:"productCode"`
	CustomerId     string         `json:"customerId"`
	ProductDetails ProductDetails `json:"productDetails"`
}

type ProductDetails struct {
	Install Install `json:"install"`
}

type Install struct {
	OperatingSystem string `json:"operatingSystem"`
}

type OrderResponse struct {
	Code    string `json:"code"`
	OrderId int    `json:"orderId"`
}

type ConfirmRequest struct {
	ApiUuid    string `json:"APIUUID"`
	ApiSecret  string `json:"APISECRET"`
	Command    string `json:"cmd"`
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
}

type ConfirmResponse struct {
	Code       string `json:"code"`
	ContractId []int  `json:"contractID"`
}
