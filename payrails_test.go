package payrails_test

import (
	"errors"
	"testing"

	payrails "github.com/toniastro/payrails-ng"
)

func setup(t *testing.T) *payrails.Registry {
	t.Helper()
	r, err := payrails.Load()
	if err != nil {
		t.Fatalf("failed to load registry: %v", err)
	}
	return r
}

func TestFindByNIBSS(t *testing.T) {
	r := setup(t)

	bank, err := r.FindByNIBSS("000013")
	if err != nil {
		t.Fatalf("expected bank, got error: %v", err)
	}
	if bank.Name != "Guaranty Trust Bank" {
		t.Errorf("expected Guaranty Trust Bank, got %s", bank.Name)
	}
}

func TestFindByNIBSS_NotFound(t *testing.T) {
	r := setup(t)

	_, err := r.FindByNIBSS("000000")
	if !errors.Is(err, payrails.ErrBankNotFound) {
		t.Errorf("expected ErrBankNotFound, got %v", err)
	}
}

func TestFindByName(t *testing.T) {
	r := setup(t)

	// "Guaranty Trust Bank" is the full name from the dataset (title-cased)
	cases := []string{"Guaranty Trust Bank", "guaranty trust bank", "GUARANTY TRUST BANK"}
	for _, name := range cases {
		bank, err := r.FindByName(name)
		if err != nil {
			t.Errorf("FindByName(%q) failed: %v", name, err)
			continue
		}
		if bank.NIBSSCode != "000013" {
			t.Errorf("FindByName(%q): expected 000013, got %s", name, bank.NIBSSCode)
		}
	}
}

func TestFindByProviderCode(t *testing.T) {
	r := setup(t)

	bank, err := r.FindByProviderCode(payrails.Paystack, "058")
	if err != nil {
		t.Fatalf("expected bank, got error: %v", err)
	}
	if bank.NIBSSCode != "000013" {
		t.Errorf("expected 000013, got %s", bank.NIBSSCode)
	}
}

func TestResolve(t *testing.T) {
	r := setup(t)

	cases := []struct {
		nibss    string
		provider payrails.Provider
		want     string
	}{
		{"000013", payrails.Paystack, "058"},    // GTBank -> Paystack
		{"000013", payrails.Nomba, "058"},       // GTBank -> Nomba
		{"000013", payrails.Mono, "058"},        // GTBank -> Mono
		{"000014", payrails.Paystack, "044"},    // Access Bank -> Paystack
		{"100004", payrails.Paystack, "999992"}, // OPay -> Paystack
		{"090405", payrails.Paystack, "50515"},  // Moniepoint -> Paystack
		{"090267", payrails.Paystack, "50211"},  // Kuda -> Paystack
	}

	for _, tc := range cases {
		code, err := r.Resolve(tc.nibss, tc.provider)
		if err != nil {
			t.Errorf("Resolve(%s, %s) error: %v", tc.nibss, tc.provider, err)
			continue
		}
		if code != tc.want {
			t.Errorf("Resolve(%s, %s): want %s, got %s", tc.nibss, tc.provider, tc.want, code)
		}
	}
}

func TestResolve_NoProviderCode(t *testing.T) {
	r := setup(t)

	_, err := r.Resolve("000013", "unknown_provider")
	if !errors.Is(err, payrails.ErrNoProviderCode) {
		t.Errorf("expected ErrNoProviderCode, got %v", err)
	}
}
