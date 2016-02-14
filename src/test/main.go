package main

import (
	//	"bytes"
	"fmt"
	"github.com/DATA-DOG/godog"
	"io/ioutil"
	"os"
	"os/exec"
	//	"strings"
)

func aBhyveInstallation() error {
	out, err := exec.Command("bhyve").CombinedOutput()
	_ = err

	if string(out[0:12]) != "Usage: bhyve" {
		return fmt.Errorf("Bhyve is not installed or not on the system PATH\n")
	}

	return nil
}

func noVmNamed(arg1 string) error {
	return godog.ErrPending
}

func iCdTo(arg1 string) error {

	err := os.Chdir("./test/config")

	if err != nil {
		return err
	}

	return nil
}

func thereIsAVmNamed(arg1 string) error {
	return godog.ErrPending
}

func aTestConfigFile() error {

	configFile, err := os.Open("./test/config/testConfig.bk")
	_ = configFile
	if os.IsNotExist(err) {
		return fmt.Errorf("testConfig.bk does not exist")
	}

	return nil
}
func bkShouldOutputonTheFirstLine(arg1 string) error {

	output, err := ioutil.ReadFile("bk.out")

	if err != nil {
		return err
	}

	if string(output) != "Config File Parsed" {
		return fmt.Errorf("Config File NOT Parsed")
	}

	return nil
}

func main() {

	//	_ = aBhyveInstallation()

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

		s.Step(`^a Bhyve installation$`, aBhyveInstallation)
		s.Step(`^no vm named: "([^"]*)"$`, noVmNamed)
		s.Step(`^I cd to: "([^"]*)"$`, iCdTo)
		s.Step(`^I run: "([^"]*)"$`, iRun)
		s.Step(`^there is a vm named: "([^"]*)"$`, thereIsAVmNamed)
		s.Step(`^a test config file$`, aTestConfigFile)
		s.Step(`^bk should output "([^"]*)" on the first line$`, bkShouldOutputonTheFirstLine)
	})
}
