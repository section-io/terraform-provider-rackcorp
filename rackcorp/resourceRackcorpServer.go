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
			"root_password": {
				Type:     schema.TypeString,
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
			"device_cancel_transaction_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"device_cancel_transaction_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRackcorpServerPopulateFromDevice(d *schema.ResourceData, config Config) error {
	deviceId := d.Get("device_id").(string)
	log.Printf("[TRACE] Rackcorp device id '%s'", deviceId)

	if deviceId == "" {
		return nil
	}

	device, err := config.Client.DeviceGet(deviceId)
	if err != nil {
		return errors.Wrapf(err, "Could not get Rackcorp device with id '%s'.", deviceId)
	}

	log.Printf("[DEBUG] Rackcorp device: %#v", device)

	panicOnError(d.Set("name", device.Name))
	panicOnError(d.Set("primary_ip", device.PrimaryIP))

	powerSwitch := getExtraByKey("SYS_POWERSWITCH", device.Extra)
	if powerSwitch == "ONLINE" {
		powerStatus := getExtraByKey("SYS_POWERSTATUS", device.Extra)
		log.Printf("[TRACE] Rackcorp device power status: %s", powerStatus)
		panicOnError(d.Set("device_status", powerStatus))
	} else {
		log.Printf("[TRACE] Rackcorp device power switch: %s", powerSwitch)
		panicOnError(d.Set("device_status", powerSwitch))
	}

	return nil
}

func resourceRackcorpServerPopulateFromContract(d *schema.ResourceData, config Config) error {
	contractId := d.Get("contract_id").(string)
	log.Printf("[TRACE] Rackcorp contract id '%s'", contractId)

	if contractId == "" {
		return nil
	}

	contract, err := config.Client.OrderContractGet(contractId)
	if err != nil {
		return errors.Wrapf(err, "Could not get Rackcorp contract with id '%s'.", contractId)
	}

	log.Printf("[DEBUG] Rackcorp contract: %#v", contract)

	panicOnError(d.Set("contract_status", contract.Status))
	panicOnError(d.Set("device_id", contract.DeviceId))

	return nil
}

func resourceRackcorpServerPopulateFromTransaction(d *schema.ResourceData, config Config) error {
	cancelTransactionId := d.Get("device_cancel_transaction_id").(int)
	log.Printf("[TRACE] Rackcorp TransactionId id '%d'", cancelTransactionId)

	if cancelTransactionId == 0 {
		return nil
	}

	transaction, err := config.Client.TransactionGet(cancelTransactionId)
	if err != nil {
		return errors.Wrapf(err, "Could not get Rackcorp transaction with id '%s'.", cancelTransactionId)
	}

	log.Printf("[DEBUG] Rackcorp transaction: %#v", transaction)

	panicOnError(d.Set("device_cancel_transaction_status", transaction.Status))

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

func startServer(deviceId string, config Config) error {
	transaction, err := config.Client.TransactionCreate(
		api.TransactionTypeStartup,
		api.TransactionObjectTypeDevice,
		deviceId,
		false)

	if err != nil {
		return errors.Wrapf(err, "Failed to start server with device id '%s'.", deviceId)
	}

	log.Printf("[TRACE] Created transaction '%d' to start server with device id '%s'.",
		transaction.TransactionId, deviceId)

	return nil
}

func cancelServer(deviceId string, d *schema.ResourceData, config Config) error {
	transaction, err := config.Client.TransactionCreate(
		api.TransactionTypeCancel,
		api.TransactionObjectTypeDevice,
		deviceId,
		true)

	panicOnError(d.Set("device_cancel_transaction_id", transaction.TransactionId))

	err = waitForTransactionAttribute(d, config, "device_cancel_transaction_status", "COMPLETED", []string{"PENDING", "COMMENCED"})

	if err != nil {
		return errors.Wrapf(err, "Failed to cancel server with device id '%s'.", deviceId)
	}

	log.Printf("[TRACE] Created transaction '%d' to cancel server with device id '%s'.",
		transaction.TransactionId, deviceId)

	return nil
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	credentials := api.Credentials{
		RootPassword: d.Get("root_password").(string),
	}

	install := api.Install{
		OperatingSystem: d.Get("operating_system").(string),
	}

	productDetails := api.ProductDetails{
		Credentials: credentials,
		Install:     install,
		CpuCount:    d.Get("cpu_count").(int),
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

	d.SetId(orderId)

	contractCount := len(confirmedOrder.ContractIds)
	if contractCount != 1 {
		return errors.Errorf("Expected one Rackcorp contract for order '%s' but received %d", orderId, contractCount)
	}

	contractId := confirmedOrder.ContractIds[0]

	panicOnError(d.Set("contract_id", contractId))

	err = waitForContractStatus(d, config, "ACTIVE", []string{"PENDING"})
	if err != nil {
		return errors.Wrapf(err, "Error waiting for Rackcorp contract status to be ACTIVE '%s'.", err)
	}

	deviceId := d.Get("device_id").(string)
	err = startServer(deviceId, config)
	if err != nil {
		return err
	}

	err = waitForDeviceAttribute(d, config, "device_status", "ONLINE", []string{"OFFLINE"})
	if err != nil {
		return errors.Wrapf(err, "Error waiting for Rackcorp device status to be ONLINE '%s'.", err)
	}

	// TODO wait for all pending transactions for device to complete

	return nil
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

	err = resourceRackcorpServerPopulateFromContract(d, config)
	if err != nil {
		return err
	}

	return resourceRackcorpServerPopulateFromDevice(d, config)
}

func resourceRackcorpServerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)
	deviceId := d.Get("device_id").(string)
	err := cancelServer(deviceId, d, config)
	if err != nil {
		return err
	}
	return nil
}

func waitForContractStatus(d *schema.ResourceData, config Config, targetStatus string, pendingStatuses []string) error {
	log.Printf(
		"[INFO] Waiting for contract to have Status of %s",
		targetStatus)

	stateConf := &resource.StateChangeConf{
		Pending:    pendingStatuses,
		Target:     []string{targetStatus},
		Refresh:    newContractStatusRefreshFunc(d, config),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func newContractStatusRefreshFunc(d *schema.ResourceData, config Config) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := resourceRackcorpServerPopulateFromContract(d, config)
		if err != nil {
			return nil, "", err
		}

		if status, ok := d.GetOk("contract_status"); ok {
			return d, status.(string), nil
		}

		return d, "", nil
	}
}

func waitForDeviceAttribute(
	d *schema.ResourceData, config Config, attribute string, target string, pending []string) error {

	log.Printf(
		"[INFO] Waiting for device to have %s of %s",
		attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newDeviceStateRefreshFunc(d, config, attribute),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func newDeviceStateRefreshFunc(d *schema.ResourceData, config Config, attribute string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		err := resourceRackcorpServerPopulateFromDevice(d, config)
		if err != nil {
			return nil, "", err
		}

		if status, ok := d.GetOk(attribute); ok {
			return d, status.(string), nil
		}

		return d, "", nil
	}
}

func waitForTransactionAttribute(
	d *schema.ResourceData, config Config, attribute string, target string, pending []string) error {

	log.Printf(
		"[INFO] Waiting for transaction to have %s of %s",
		attribute, target)

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newTransactionStateRefreshFunc(d, config, attribute),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func newTransactionStateRefreshFunc(d *schema.ResourceData, config Config, attribute string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		err := resourceRackcorpServerPopulateFromTransaction(d, config)
		if err != nil {
			return nil, "", err
		}

		if status, ok := d.GetOk(attribute); ok {
			return d, status.(string), nil
		}

		return d, "", nil
	}
}
