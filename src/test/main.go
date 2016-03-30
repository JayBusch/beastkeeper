package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
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

	if _, err := os.Stat("test/virtualMachines/" + arg1 + ".img"); os.IsNotExist(err) {
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

	_, fileErr := os.Stat("test/config/testConfig.bk")

	if os.IsNotExist(fileErr) {
		return fileErr
	}

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

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func runningUnameOverSSHOnTheInstanceYields(arg1 string) error {

	return godog.ErrPending

	sshConfig := &ssh.ClientConfig{
		User: "bk",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	connection, err := ssh.Dial("tcp", "127.0.0.1:10001", sshConfig)
	if err != nil {
		fmt.Errorf("Failed to dial: %s", err)
	}

	session, sessErr := connection.NewSession()

	if err != nil {
		fmt.Errorf("Failed to create session: %s", err)
	}

	_ = session
	_ = sessErr

	return godog.ErrPending
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
		s.Step(`^running uname over SSH on the instance yields "([^"]*)"$`,
			runningUnameOverSSHOnTheInstanceYields)
	})
}
