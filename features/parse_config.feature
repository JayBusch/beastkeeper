Feature: parse a config file

	Scenario: parse test config file
	  Given a test config file
	  When I cd to: "test"
	  And I run: "../bin/bk -c ./config/testConfig.bk config print > bk.out"
	  Then bk output should match testConfig
