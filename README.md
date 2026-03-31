# payrails-ng

A Go library for resolving Nigerian bank codes across multiple payment providers. It maps NIBSS (Nigerian Inter-Bank Settlement System) codes to provider-specific bank codes, making it easy to integrate with various payment rails.

## Features

- Lookup banks by NIBSS code, name/alias, or provider-specific code
- Resolve NIBSS codes to provider-specific codes for payouts
- Case-insensitive, whitespace-tolerant name matching
- 172 Nigerian banks and financial institutions
- Zero external dependencies — bank data is embedded at compile time

## Supported Providers

`paystack` · `flutterwave` · `monnify` · `nomba` · `redbiller` · `mono` · `aella` · `fincra` · `budpay` · `payaza`

## Installation

```bash
go get github.com/toniastro/payrails-ng
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	payrails "github.com/toniastro/payrails-ng"
)

func main() {
	registry, err := payrails.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Resolve a NIBSS code to a provider-specific code
	code, err := registry.Resolve("000001", payrails.Paystack)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code) // "232"

	// Find a bank by name
	bank, err := registry.FindByName("Sterling Bank")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bank.NIBSSCode) // "000001"

	// Find a bank by NIBSS code
	bank, err = registry.FindByNIBSS("000013")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bank.Name) // "GT Bank"

	// Find a bank by provider-specific code
	bank, err = registry.FindByProviderCode(payrails.Paystack, "232")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bank.Name) // "Sterling Bank"

	// List all banks
	banks := registry.All()
	fmt.Printf("%d banks loaded\n", len(banks))
}
```

## Contributing

Contributions are welcome! If a bank is missing or a provider code is incorrect, update `data/banks.json` and open a pull request.

## License

MIT