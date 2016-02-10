Feature: build a vm

	Scenario: build FreeBSD_10_2_base vm
	  Given a Bhyve installation
	  And no vm named: "test_FreeBSD_10_2_base"
	  When I cd to: "virtmachines/"
	  And I run: "bk newvm -v 10.2"
	  Then there is a vm named: "test_FreeBSD_10_2_base"

