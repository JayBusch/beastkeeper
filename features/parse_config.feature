Feature: parse a config file

	Scenario: parse test config file
	  Given a test config file
	  When I cd to: "test/config"
	  And I run: "../../bin/bk -c testConfig.bk config print > bk.out"
	  Then bk should output "Config File Parsed As:" on the first line

