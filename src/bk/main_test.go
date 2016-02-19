package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// Test config JSON, we are going to write this to the
// file system and allow the parseConfigFile function to read it back.
// The idea being that unit tests should be self contained and not rely
// on external resources.
var testConfigFileName = "../../test/config/unitTest_config_file.bk"

var testConfigAsBytes = []byte(`
	{
		"Instances": [{
			"ID" : "2237b51a-39aa-4720-b68a-d1ee214e9272",
			"Label" : "test_instance_1",
			"Type" : "VM",
			"Address" : "192.242.2.200",
			"AdminLogin" : "root",
			"Containers" : []
			}]	
	}
	`)

func writeTestConfigToFile() {
	ioutil.WriteFile(testConfigFileName, testConfigAsBytes, 0644)
}

func TestParseConfigFile(t *testing.T) {

	writeTestConfigToFile()

	testConfig := parseConfigFile(testConfigFileName)

	// Check the values on our returned config file to see if they match
	if len(testConfig.Instances) < 1 {
		t.Fatalf("No data was unmarshalled")
	} else if testConfig.Instances[0].ID.UUID.String() != "2237b51a-39aa-4720-b68a-d1ee214e9272" {
		t.Fatalf("ID does not match")

	} else if testConfig.Instances[0].Label != "test_instance_1" {
		t.Fatalf("Label does not match")
	} else if testConfig.Instances[0].Type != VM {
		t.Fatalf("Type does not match")
	} else if testConfig.Instances[0].Address.String() != "192.242.2.200" {
		t.Fatalf("Address does not match")
	} else if testConfig.Instances[0].AdminLogin != "root" {
		t.Fatalf("AdminLogin does not match")
	}
}

func TestCommandConfigPrint(t *testing.T) {

	writeTestConfigToFile()
	parseConfigFile(testConfigFileName)

	old_stdout := os.Stdout
	r, w, err := os.Pipe()
	os.Stdout = w

	if err != nil {
		t.Fatalf("Could not open Pipe to STDOUT")
	}

	commandConfigPrint()

	outC := make(chan []byte)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.Bytes()
	}()

	w.Close()
	os.Stdout = old_stdout
	out := <-outC

	buffer := new(bytes.Buffer)

	if jsonErr := json.Compact(buffer, testConfigAsBytes); jsonErr != nil {
		t.Fatalf("Error compacting testConfigAsBytes")
	}

	if !bytes.Equal(out, buffer.Bytes()) {
		t.Logf("out: %v", string(out))
		t.Logf("buffer: %v", buffer)
		t.Fatalf("Config file not printed correctly")
	}
}
