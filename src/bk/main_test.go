package main

import (
	"io/ioutil"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	// Test config as JSON string, we are going to write this to the
	// file system and allow the parseConfigFile function to read it back.
	// The idea being that unit tests should be self contained and not rely
	// on external resources.
	testConfigAsString := `
	{
		"instances": [{
			"ID" : "2237b51a-39aa-4720-b68a-d1ee214e9272",
			"Label" : "test_instance_1",
			"Type" : "VM",
			"Address" : "192.242.2.200"
			}]	
	}
	`

	testConfigFileName := "../../test/config/unitTest_config_file.bk"

	ioutil.WriteFile(testConfigFileName, []byte(testConfigAsString), 0644)

	testConfig := parseConfigFile(testConfigFileName)

	// Check the values on our returned config file to see if they match
	if len(testConfig.Instances) < 1 {
		t.Fatalf("No data was unmarshalled")
	} else if testConfig.Instances[0].Label != "test_instance_1" {
		t.Fatalf("Label does not match")
	} else if testConfig.Instances[0].Type != VM {
		t.Fatalf("Type does not match")
	} else if testConfig.Instances[0].Address.String() != "192.242.2.200" {
		t.Fatalf("Address does not match")
	}
}
