package main

import (
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	test := parseConfigFile("testconfig.bk")
	_ = test
	t.Fatalf("Initial Failing Test")
}
