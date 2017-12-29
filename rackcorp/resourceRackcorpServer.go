package rackcorp

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/section-io/terraform-provider-rackcorp/rackcorp/api"
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
			"contract_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contract_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getDeviceByContract(contractId string, d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	panicOnError(d.Set("contract_id", contractId))

	contract, err := waitForContractStatus(contractId, "ACTIVE", []string{"PENDING"}, config.Client)

	if err != nil {
		return errors.Wrapf(err, "Error waiting for Rackcorp contract status to be ACTIVE '%s'.", err)
	}

	panicOnError(d.Set("contract_status", contract.Status))

	deviceId := contract.DeviceId
	if deviceId == "" {
		log.Printf("[WARN] Rackcorp contract '%s' device ID not specified.", contractId)
		d.SetId("")
		return nil
	}

	device, err := config.Client.DeviceGet(deviceId)
	if err != nil {
		return errors.Wrapf(err, "Error retrieving Rackcorp device '%s'.", deviceId)
	}

	if device.DeviceId == "" {
		log.Printf("[WARN] Rackcorp device '%s' not found.", deviceId)
		d.SetId("")
		return nil
	}

	sysPowerSwitch := getExtraByKey("SYS_POWERSWITCH", device.Extra)
	sysPowerStatus := getExtraByKey("SYS_POWERSTATUS", device.Extra)

	if sysPowerSwitch == "ONLINE" && sysPowerStatus == "ONLINE" {
		panicOnError(d.Set("device_status", "ONLINE"))
	}

	_, err = waitForDeviceAttribute(d, "ONLINE", []string{""}, "device_status", meta)

	if err != nil {
		return errors.Wrapf(err, "Error waiting for Rackcorp device status to be ONLINE '%s'.", err)
	}

	log.Printf("[DEBUG] Rackcorp device: %#v", device)

	panicOnError(d.Set("device_id", device.DeviceId))
	panicOnError(d.Set("name", device.Name))
	panicOnError(d.Set("primary_ip", device.PrimaryIP))

	return nil
}

func getExtraByKey(key string, extras []api.DeviceExtra) string {
	for _, extra := range extras {
		if extra.Key == key {
			if extra.Value == nil {
				return ""
			}
			if s, ok := extra.Value.(string); ok {
				return s
			}
			return ""
		}
	}
	return ""
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	install := api.Install{
		OperatingSystem: d.Get("operating_system").(string),
	}

	productDetails := api.ProductDetails{
		Install:  install,
		CpuCount: d.Get("cpu_count").(int),
	}

	productCode := "SERVER_VIRTUAL_" + d.Get("server_class").(string) + "_" + d.Get("country").(string)

	createdOrder, err := config.Client.OrderCreate(productCode, config.CustomerId, productDetails)
	if err != nil {
		return errors.Wrap(err, "Rackcorp order create request failed.")
	}

	orderId := createdOrder.OrderId
	confirmedOrder, err := config.Client.OrderConfirm(orderId)
	if err != nil {
		return errors.Wrapf(err, "Failed to confirm Rackcorp server order '%s'.", orderId)
	}

	contractCount := len(confirmedOrder.ContractIds)
	if contractCount != 1 {
		return errors.Errorf("Expected one Rackcorp contract for order '%s' but received %d", orderId, contractCount)
	}

	contractId := confirmedOrder.ContractIds[0]

	panicOnError(d.Set("contract_id", contractId))

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

	order, err := config.Client.OrderGet(orderId)
	if err != nil {
		return errors.Wrapf(err, "Error retrieving Rackcorp order '%s'.", orderId)
	}

	contractId := order.ContractId
	if contractId == "" {
		log.Printf("[WARN] Rackcorp order '%s' not found.", orderId)
		d.SetId("")
		return nil
	}
	panicOnError(d.Set("contract_id", contractId))

	return getDeviceByContract(contractId, d, config)
}

func resourceRackcorpServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func waitForContractStatus(contractId string, targetStatus string, pendingStatuses []string, client api.Client) (*api.OrderContract, error) {
	log.Printf(
		"[INFO] Waiting for contract (%s) to have Status of %s",
		contractId, targetStatus)

	stateConf := &resource.StateChangeConf{
		Pending:    pendingStatuses,
		Target:     []string{targetStatus},
		Refresh:    newContractStateRefreshFunc(contractId, client),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	result, err := stateConf.WaitForState()
	if err != nil {
		return nil, err
	}
	contract := result.(*api.OrderContract)
	return contract, nil
}

func newContractStateRefreshFunc(contractId string, client api.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		contract, err := client.OrderContractGet(contractId)
		if err != nil {
			return nil, "", errors.Errorf("Error retrieving contract: %v", err)
		}

		return contract, contract.Status, nil
	}
}

func waitForDeviceAttribute(
	d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	// Wait for the contract so we can get the device attributes
	// that show up after a while

	deviceId := d.Get("device_id").(string)
	log.Printf(
		"[INFO] Waiting for device (%s) to have %s of %s",
		deviceId, attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newDeviceStateRefreshFunc(deviceId, meta.(Config).Client),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForState()
}

func newDeviceStateRefreshFunc(deviceId string, client api.Client) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {

		device, err := client.DeviceGet(deviceId)
		if err != nil {
			return nil, "", errors.Wrapf(err, "Error retrieving device id '%s'", deviceId)
		}

		return device, device.Status, nil
	}
}
