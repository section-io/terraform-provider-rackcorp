package rackcorp

import (
	"log"
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
			"cpu_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getDeviceByContract(contractId string, d *schema.ResourceData, config Config) error {
	contractGetRequest := NewOrderContractGetRequest(contractId)
	contractGetResponse, err := contractGetRequest.Post(config)
	if err != nil {
		return errors.Wrapf(err, "Error retrieving Rackcorp contract '%s'.", contractId)
	}

	log.Printf("[DEBUG] Rackcorp get contract response: %#v", contractGetResponse)

	if contractGetResponse.Contract.ContractId == "" {
		log.Printf("[WARN] Rackcorp contract '%s' not found.", contractId)
		d.SetId("")
		return nil
	}

	if contractGetResponse.Contract.Status == RackcorpApiOrderContractStatusPending {
		log.Printf("[WARN] Rackcorp contract '%s' is pending.", contractId)
		return nil
		// TODO implement waiting with retry, eg:
		//  https://github.com/terraform-providers/terraform-provider-digitalocean/blob/master/digitalocean/resource_digitalocean_droplet.go#L562
	}

	deviceId := contractGetResponse.Contract.DeviceId
	if deviceId == "" {
		log.Printf("[WARN] Rackcorp contract '%s' device ID not specified.", contractId)
		d.SetId("")
		return nil
	}

	deviceGetRequest := NewDeviceGetRequest(deviceId)
	deviceGetResponse, err := deviceGetRequest.Post(config)
	if err != nil {
		return errors.Wrapf(err, "Error retrieving Rackcorp device '%s'.", deviceId)
	}

	device := deviceGetResponse.Device
	if device.DeviceId == "" {
		log.Printf("[WARN] Rackcorp device '%s' not found.", deviceId)
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Rackcorp device: %#v", device)

	panicOnError(d.Set("device_id", device.Name))
	panicOnError(d.Set("name", device.Name))
	panicOnError(d.Set("primary_ip", device.PrimaryIP))

	return nil
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	install := Install{
		OperatingSystem: d.Get("operating_system").(string),
	}

	productDetails := ProductDetails{
		Install:  install,
		CpuCount: d.Get("cpu_count").(int),
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

	orderId := strconv.Itoa(orderResponse.OrderId)
	confirmRequest := NewOrderConfirmRequest(orderId)

	confirmResponse, err := confirmRequest.Post(config)
	if err != nil {
		return errors.Wrapf(err, "Failed to confirm Rackcorp server order '%s'.", orderId)
	}

	if confirmResponse.Code != RackcorpApiResponseCodeOK {
		return errors.Errorf("Unexpected Rackcorp server order response code '%s'.", confirmResponse.Code)
		// TODO log message too
	}

	contractCount := len(confirmResponse.ContractIds)
	if contractCount != 1 {
		return errors.Errorf("Expected one Rackcorp contract for order '%s' but received %d", orderId, contractCount)
	}

	contractId := strconv.Itoa(confirmResponse.ContractIds[0])

	d.SetId(orderId)

	return getDeviceByContract(contractId, d, config)
}

func panicOnError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func resourceRackcorpServerRead(d *schema.ResourceData, meta interface{}) error {
	orderId := d.Id()
	if orderId == "" {
		return errors.Errorf("Missing resource id.")
	}

	config := meta.(Config)

	orderGetRequest := NewOrderGetRequest(orderId)
	orderGetResponse, err := orderGetRequest.Post(config)
	if err != nil {
		return errors.Wrapf(err, "Error retrieving Rackcorp order '%s'.", orderId)
	}

	contractId := orderGetResponse.Order.ContractId
	if contractId == "" {
		log.Printf("[WARN] Rackcorp order '%s' not found.", orderId)
		d.SetId("")
		return nil
	}

	return getDeviceByContract(contractId, d, config)
}

func resourceRackcorpServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
