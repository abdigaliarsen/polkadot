package main

import (
	"strings"
)

var kusamaPrefixLetters = "CDEFGHJ"

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

// IsPolkadot https://wiki.polkadot.network/learn/learn-account-advanced/#address-format
func IsPolkadot(address string) bool {
	return strings.HasPrefix(address, "1")
}

// IsKusama https://wiki.polkadot.network/learn/learn-account-advanced/#address-format
func IsKusama(address string) bool {
	for _, l := range kusamaPrefixLetters {
		if strings.HasPrefix(address, string(l)) {
			return true
		}
	}

	return false
}

// IsGenericSubstrate https://wiki.polkadot.network/learn/learn-account-advanced/#address-format
func IsGenericSubstrate(address string) bool {
	return strings.HasPrefix(address, "5")
}
