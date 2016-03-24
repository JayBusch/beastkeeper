Feature: Enforce Test Configuration 

	Scenario: enforce the config in test/config/testConfig.bk
	  Given a Bhyve installation
	  And no vm named: "test_instance_1"
	  When I cd to: "test"
	  And I run: "../bin/bk -c ./config/testConfig.bk enforce"
	  Then there is a vm named: "test_instance_1"
	  And running uname over SSH on the instance yields "FreeBSD"
	  
