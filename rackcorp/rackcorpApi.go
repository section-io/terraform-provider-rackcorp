package rackcorp

// TODO move these to enums in the api package
const (
	RackcorpApiResponseCodeOK           = "OK"
	RackcorpApiResponseCodeAccessDenied = "ACCESS_DENIED"
	RackcorpApiResponseCodeFault        = "FAULT"

	RackcorpApiOrderStatusPending  = "PENDING"
	RackcorpApiOrderStatusAccepted = "ACCEPTED"

	RackcorpApiOrderContractStatusActive  = "ACTIVE"
	RackcorpApiOrderContractStatusPending = "PENDING"
)
