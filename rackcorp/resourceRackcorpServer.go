package rackcorp

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

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
		OperatingSystem: d.Get("operating_system").(string),
	}

	productDetails := ProductDetails{
		Install: install,
	}

	orderRequest := NewOrderCreateRequest()
	orderRequest.CustomerId = config.CustomerId
	orderRequest.ProductCode = "SERVER_VIRTUAL_" + d.Get("server_class").(string) + "_" + d.Get("country").(string)
	orderRequest.ProductDetails = productDetails

	orderResponse, err := orderRequest.Post(config)
	if err != nil {
		return errors.Wrap(err, "Rackcorp order create request failed.")
	}

	if orderResponse.Code != RackcorpApiResponseCodeOK {
		return errors.Errorf("Unexpected Rackcorp server order response code '%s'.", orderResponse.Code)
		// TODO log message too
	}

	confirmRequest := NewOrderConfirmRequest(strconv.Itoa(orderResponse.OrderId))

	confirmResponse, err := confirmRequest.Post(config)
	if err != nil {
		return errors.Wrapf(err, "Failed to confirm Rackcorp server order '%s'.", confirmRequest.OrderId)
	}

	if confirmResponse.Code != RackcorpApiResponseCodeOK {
		return errors.Errorf("Unexpected Rackcorp server order response code '%s'.", confirmResponse.Code)
		// TODO log message too
	}

	return resourceRackcorpServerRead(d, meta)
}

func resourceRackcorpServerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRackcorpServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
