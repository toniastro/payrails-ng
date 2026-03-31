package payrails

import (
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
)

//go:embed data/banks.json
var banksData []byte

var (
	ErrBankNotFound   = errors.New("bank not found")
	ErrNoProviderCode = errors.New("no code found for provider")
)

// Registry holds all banks and their lookup indexes
type Registry struct {
	banks      []*Bank
	byNIBSS    map[string]*Bank
	byAlias    map[string]*Bank
	byProvider map[Provider]map[string]*Bank
}

// Load initializes the registry from the embedded banks.json
func Load() (*Registry, error) {
	var banks []*Bank
	if err := json.Unmarshal(banksData, &banks); err != nil {
		return nil, err
	}

	r := &Registry{
		banks:      banks,
		byNIBSS:    make(map[string]*Bank),
		byAlias:    make(map[string]*Bank),
		byProvider: make(map[Provider]map[string]*Bank),
	}

	for _, b := range banks {
		// Index by NIBSS code
		r.byNIBSS[b.NIBSSCode] = b

		// Index by name and aliases (lowercased for case-insensitive lookup)
		r.byAlias[strings.ToLower(b.Name)] = b
		r.byAlias[strings.ToLower(b.ShortName)] = b
		for _, alias := range b.Aliases {
			r.byAlias[strings.ToLower(alias)] = b
		}

		// Index by provider code
		for provider, code := range b.ProviderCodes {
			if r.byProvider[provider] == nil {
				r.byProvider[provider] = make(map[string]*Bank)
			}
			r.byProvider[provider][code] = b
		}
	}

	return r, nil
}

// FindByNIBSS looks up a bank by its NIBSS code
func (r *Registry) FindByNIBSS(nibssCode string) (*Bank, error) {
	b, ok := r.byNIBSS[nibssCode]
	if !ok {
		return nil, ErrBankNotFound
	}
	return b, nil
}

// FindByName looks up a bank by name or alias (case-insensitive)
func (r *Registry) FindByName(name string) (*Bank, error) {
	b, ok := r.byAlias[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return nil, ErrBankNotFound
	}
	return b, nil
}

// FindByProviderCode looks up a bank by a provider-specific code
func (r *Registry) FindByProviderCode(p Provider, code string) (*Bank, error) {
	providerIndex, ok := r.byProvider[p]
	if !ok {
		return nil, ErrBankNotFound
	}
	b, ok := providerIndex[code]
	if !ok {
		return nil, ErrBankNotFound
	}
	return b, nil
}

// Resolve returns the provider-specific code for a bank given its NIBSS code.
// This is the primary API for payout rail resolution.
func (r *Registry) Resolve(nibssCode string, p Provider) (string, error) {
	b, err := r.FindByNIBSS(nibssCode)
	if err != nil {
		return "", err
	}
	code, ok := b.ProviderCodes[p]
	if !ok {
		return "", ErrNoProviderCode
	}
	return code, nil
}

// All returns a copy of all banks in the registry
func (r *Registry) All() []*Bank {
	out := make([]*Bank, len(r.banks))
	copy(out, r.banks)
	return out
}
