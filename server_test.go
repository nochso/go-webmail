package main

import (
	"net/mail"
	"testing"
)

var belongsTests = []struct {
	expected bool
	addr     string
}{
	{true, "user@domain.tld"},
	{false, "missing@nope.tld"},
	{false, "user@nope.domain.tld"},
	{false, "user.domain.tld@nope.tld"},
	{true, "user@some.subdomain.tld"},
	{true, "user@some.more.subdomain.tld"},
	{false, "user.some.subdomain.tld@nope.tld"},
}

func TestAddressBelongsToHost(t *testing.T) {
	cfg.Domains = []string{
		"domain.tld",
		".subdomain.tld",
	}
	for _, tt := range belongsTests {
		addr := &mail.Address{Address: tt.addr}
		if addressBelongsToHost(addr) != tt.expected {
			t.Fail()
		}
	}

	cfg.Domains = []string{"*"}
	for _, tt := range belongsTests {
		addr := &mail.Address{Address: tt.addr}
		if !addressBelongsToHost(addr) {
			t.Errorf("Address '%s' must be accepted by accepted by catch-all '*' configuration via smtp.accept.domains", addr)
		}
	}
}
