package rackcorp

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pkg/errors"
	"github.com/section-io/rackcorp-sdk-go/api"
)

func storageSchemaElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					api.StorageTypes,
					false,
				),
			},
			"sort_order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func firewallPolicySchemaElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					api.FirewallPolicyDirections,
					false,
				),
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					api.FirewallPolicyTypes,
					false,
				),
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func nicSchemaElement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"speed": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ipv4": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"pool_ipv4": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pool_ipv6": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
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
				ValidateFunc: validation.StringInSlice(
					api.ServerClasses,
					false,
				),
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
			"memory_gb": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"root_password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"data_center_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"traffic_gb": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"post_install_script": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     storageSchemaElement(),
			},
			"firewall_policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     firewallPolicySchemaElement(),
			},
			"nics": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     nicSchemaElement(),
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
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
		if apiErr, ok := err.(*api.ApiError); ok {
			if apiErr.Message == "Could not find device" {
				return newNotFoundError(apiErr.Message)
			}
		}
		return errors.Wrapf(err, "Could not get Rackcorp device with id '%s'.", deviceId)
	}

	log.Printf("[DEBUG] Rackcorp device: %#v", device)

	panicOnError(d.Set("name", device.Name))
	panicOnError(d.Set("primary_ip", device.PrimaryIP))
	panicOnError(d.Set("data_center_id", device.DataCenterId))

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
	cancelTransactionId := d.Get("device_cancel_transaction_id").(string)
	log.Printf("[TRACE] Rackcorp TransactionId id '%s'.", cancelTransactionId)

	if cancelTransactionId == "" {
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

	log.Printf("[TRACE] Created transaction '%s' to start server with device id '%s'.",
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

	log.Printf("[TRACE] Created transaction '%s' to cancel server with device id '%s'.",
		transaction.TransactionId, deviceId)

	return nil
}

func translateFirewallPolicy(d *schema.ResourceData) []api.FirewallPolicy {
	var result []api.FirewallPolicy
	list, ok := d.GetOk("firewall_policies")
	if !ok {
		return result
	}

	for _, raw := range list.([]interface{}) {
		data := raw.(map[string]interface{})

		policy := api.FirewallPolicy{
			Direction: data["direction"].(string),
			Policy:    data["policy"].(string),
			Order:     data["order"].(int),
		}

		if v := data["comment"].(string); v != "" {
			policy.Comment = v
		}

		if v := data["ip_address_from"].(string); v != "" {
			policy.IpAddressFrom = v
		}

		if v := data["ip_address_to"].(string); v != "" {
			policy.IpAddressTo = v
		}

		if v := data["port_from"].(string); v != "" {
			policy.PortFrom = v
		}

		if v := data["port_to"].(string); v != "" {
			policy.PortTo = v
		}

		if v := data["protocol"].(string); v != "" {
			policy.Protocol = v
		}
		result = append(result, policy)
	}

	return result
}

func translateStorage(d *schema.ResourceData) []api.Storage {
	var result []api.Storage
	list, ok := d.GetOk("storage")
	if !ok {
		return result
	}

	for _, raw := range list.([]interface{}) {
		data := raw.(map[string]interface{})

		storage := api.Storage{
			SizeGB:      data["size_gb"].(int),
			StorageType: api.StorageTypeMagnetic,
		}

		if v := data["name"].(string); v != "" {
			storage.Name = v
		}

		if v := data["type"].(string); v != "" {
			storage.StorageType = v
		}

		if v := data["sort_order"].(int); v != 0 {
			storage.SortOrder = v
		}

		result = append(result, storage)
	}

	return result
}

func translateNic(d *schema.ResourceData) []api.Nic {
	var result []api.Nic
	list, ok := d.GetOk("nics")
	if !ok {
		return result
	}

	for _, raw := range list.([]interface{}) {
		data := raw.(map[string]interface{})

		nic := api.Nic{
			Speed: data["speed"].(int),
			IPV4:  data["ipv4"].(int),
		}

		if v := data["name"].(string); v != "" {
			nic.Name = v
		}

		if v := data["vlan"].(int); v != 0 {
			nic.Vlan = v
		}

		if v := data["pool_ipv4"].(int); v != 0 {
			nic.PoolIPv4 = v
		}

		if v := data["ipv6"].(int); v != 0 {
			nic.IPV6 = v
		}

		if v := data["pool_ipv6"].(int); v != 0 {
			nic.PoolIPv6 = v
		}

		result = append(result, nic)
	}

	return result
}

func resourceRackcorpServerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)

	credentials := []api.Credential{
		{
			Username: "root",
			Password: d.Get("root_password").(string),
		},
	}

	install := api.Install{
		OperatingSystem: d.Get("operating_system").(string),
	}

	if script, ok := d.GetOk("post_install_script"); ok {
		install.PostInstallScript = script.(string)
	}

	productDetails := api.ProductDetails{
		Credentials:      credentials,
		Install:          install,
		CpuCount:         d.Get("cpu_count").(int),
		MemoryGB:         d.Get("memory_gb").(int),
		Storage:          translateStorage(d),
		FirewallPolicies: translateFirewallPolicy(d),
		Nics:             translateNic(d),
	}

	if name, ok := d.GetOk("name"); ok {
		productDetails.Hostname = name.(string)
	}

	if dataCenterId, ok := d.GetOk("data_center_id"); ok {
		productDetails.DataCenterId = dataCenterId.(string)
	}

	if trafficGB, ok := d.GetOk("traffic_gb"); ok {
		productDetails.TrafficGB = trafficGB.(int)
	}

	productCode := api.GetVirtualServerProductCode(
		d.Get("server_class").(string),
		d.Get("country").(string),
	)

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
		return errors.Wrap(err, "Error waiting for Rackcorp contract status to be ACTIVE")
	}

	deviceId := d.Get("device_id").(string)
	err = waitForPendingDeviceTransactions(deviceId, config)
	if err != nil {
		return errors.Wrap(err, "Error waiting for Rackcorp device transactions to complete")
	}

	err = startServer(deviceId, config)
	if err != nil {
		return err
	}

	err = waitForDeviceAttribute(d, config, "device_status", "ONLINE", []string{"OFFLINE"})
	if err != nil {
		return errors.Wrap(err, "Error waiting for Rackcorp device status to be ONLINE")
	}

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

	err = resourceRackcorpServerPopulateFromDevice(d, config)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	return nil
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

func waitForPendingDeviceTransactions(deviceId string, config Config) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{api.TransactionStatusPending},
		Target:     []string{api.TransactionStatusCompleted},
		Refresh:    newPendingTransactionsRefreshFunc(deviceId, config),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func newPendingTransactionsRefreshFunc(deviceId string, config Config) resource.StateRefreshFunc {
	var dummyResource struct{}
	filter := api.TransactionFilter{
		ObjectType:   api.TransactionObjectTypeDevice,
		ObjectId:     []string{deviceId},
		Status:       []string{api.TransactionStatusPending, api.TransactionStatusCommenced},
		ResultWindow: 1,
	}

	return func() (interface{}, string, error) {

		transactions, matches, err := config.Client.TransactionGetAll(filter)
		if err != nil {
			return nil, "", err
		}

		if matches == 0 {
			return dummyResource, api.TransactionStatusCompleted, nil
		}

		for _, t := range transactions {
			log.Printf("[TRACE] pending transaction: %#v", t)
		}

		return dummyResource, api.TransactionStatusPending, nil
	}

}
