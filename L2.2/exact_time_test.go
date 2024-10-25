package main

import (
	"github.com/beevik/ntp"
	"testing"
)

func TestTime(t *testing.T) {
	address := "0.beevik-ntp.pool.ntp.org"

	_, err := ntp.Time(address)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
