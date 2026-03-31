package payrails

// Provider represents a payment provider/rail
type Provider string

const (
	Paystack    Provider = "paystack"
	Flutterwave Provider = "flutterwave"
	Monnify     Provider = "monnify"
	Nomba       Provider = "nomba"
	Redbiller   Provider = "redbiller"
	Mono        Provider = "mono"
	Aella       Provider = "aella"
	Fincra      Provider = "fincra"
	Budpay      Provider = "budpay"
	Payaza      Provider = "payaza"
)

// BankType represents the category of financial institution
type BankType string

// BankStatus represents the operational status of a bank
type BankStatus string

// Bank represents a Nigerian financial institution with its
// NIBSS code as the canonical identifier and provider-specific codes
type Bank struct {
	NIBSSCode     string              `json:"nibss_code"`
	Name          string              `json:"name"`
	ShortName     string              `json:"short_name"`
	Aliases       []string            `json:"aliases,omitempty"`
	Type          BankType            `json:"type"`
	Status        BankStatus          `json:"status"`
	ProviderCodes map[Provider]string `json:"provider_codes"`
}
