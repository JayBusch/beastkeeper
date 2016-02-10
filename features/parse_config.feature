Feature: parse a config file

	Scenario: parse test config file
	  Given a test config file
	  When I cd to: "test/config"
	  And I run: "bk -c test_config.bk"
	  Then bk should output "Config File Parsed"

