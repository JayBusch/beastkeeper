package main

import (
	//	"fmt"
	"os"

	"github.com/DATA-DOG/godog"
)

func aBhyveInstallation() error {
	return godog.ErrPending
}

func noVmNamed(arg1 string) error {
	return godog.ErrPending
}

func iCdTo(arg1 string) error {
	return godog.ErrPending
}

func thereIsAVmNamed(arg1 string) error {
	return godog.ErrPending
}

func aTestConfigFile() error {
	return godog.ErrPending
}

func bkShouldOutput(arg1 string) error {
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

		s.Step(`^a Bhyve installation$`, aBhyveInstallation)
		s.Step(`^no vm named: "([^"]*)"$`, noVmNamed)
		s.Step(`^I cd to: "([^"]*)"$`, iCdTo)
		s.Step(`^I run: "([^"]*)"$`, iRun)
		s.Step(`^there is a vm named: "([^"]*)"$`, thereIsAVmNamed)
		s.Step(`^a test config file$`, aTestConfigFile)
		s.Step(`^bk should output "([^"]*)"$`, bkShouldOutput)
	})
}
