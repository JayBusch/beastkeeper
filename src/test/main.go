package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"io/ioutil"
	"os"
	"os/exec"
)

func aBhyveInstallation() error {
	out, err := exec.Command("bhyve").CombinedOutput()
	_ = err

	if string(out[0:12]) != "Usage: bhyve" {
		return fmt.Errorf("Bhyve is not installed or not on the system PATH\n")
	}

	return nil
}

func iCdTo(arg1 string) error {

	err := os.Chdir(arg1)

	if err != nil {
		return err
	}

	return nil
}

func noVmNamed(arg1 string) error {

	//wd, _ := os.Getwd()

	//fmt.Printf("noVmNamed Running in: %v\n", wd)

	if _, err := os.Stat("test/virtualMachines/" + arg1 + ".img"); os.IsNotExist(err) {
		//fmt.Printf("Disk Image Not Found\n")
		return nil
	}

	return fmt.Errorf("%v DOES exist (and it should not)", arg1)
}

func thereIsAVmNamed(arg1 string) error {
	if _, err := os.Stat("virtualMachines/" + arg1 + ".img"); os.IsNotExist(err) {
		return err
	}

	return nil
}

func aTestConfigFile() error {

	//var workingDir, _ = os.Getwd()

	//return fmt.Errorf("%v", workingDir)

	_, fileErr := os.Stat("test/config/testConfig.bk")

	if os.IsNotExist(fileErr) {
		return fileErr
	}

	/*configFile, err := os.Open("testConfig.bk")
	_ = configFile
	if os.IsNotExist(err) {
		return fmt.Errorf("testConfig.bk does not exist")
	}*/

	return nil
}

func bkOutputShouldMatchTestConfig() error {

	outputData, err := ioutil.ReadFile("bk.out")

	if err != nil {
		return err
	}

	configFileData, configFileErr := ioutil.ReadFile("./config/testConfig.bk")

	if os.IsNotExist(configFileErr) {
		return fmt.Errorf("testConfig.bk does not exist")
	} else if configFileErr != nil {
		return configFileErr
	}

	outputDataBuffer := new(bytes.Buffer)

	odJsonErr := json.Compact(outputDataBuffer, outputData)
	if odJsonErr != nil {
		return fmt.Errorf("Error compacting outputData")
	}

	configDataBuffer := new(bytes.Buffer)

	cfgJsonErr := json.Compact(configDataBuffer, configFileData)
	if cfgJsonErr != nil {
		return fmt.Errorf("Error compacting configFileData")
	}

	if !bytes.Equal(outputDataBuffer.Bytes(), configDataBuffer.Bytes()) {
		return fmt.Errorf("Output JSON does not match original")
	}

	return nil
}

func main() {

	godog.Run(func(s *godog.Suite) {
		origWd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		s.BeforeScenario(func(_ interface{}) {
			if err := os.Chdir(origWd); err != nil {
				panic(err)
			}
		})

		s.AfterScenario(func(_ interface{}, _ error) {
			if err := os.Chdir(origWd); err != nil {
				panic(err)
			}
		})

		s.Step(`^a Bhyve installation$`, aBhyveInstallation)
		s.Step(`^no vm named: "([^"]*)"$`, noVmNamed)
		s.Step(`^I cd to: "([^"]*)"$`, iCdTo)
		s.Step(`^I run: "([^"]*)"$`, iRun)
		s.Step(`^there is a vm named: "([^"]*)"$`, thereIsAVmNamed)
		s.Step(`^a test config file$`, aTestConfigFile)
		s.Step(`^bk output should match testConfig$`, bkOutputShouldMatchTestConfig)
	})

}
