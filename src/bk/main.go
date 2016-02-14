package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	command_config_flag       = kingpin.Command("config", "work with the config file")
	command_config_print_flag = command_config_flag.Command("print", "print the config file")
	config_filename_flag      = kingpin.Flag("configFileName", "set the config file name").Short('c').String()
)

func parseConfigFile() {
	//configFileName := flag.String("c", "config.bk", "config file to load")

	configFileName := "config.bk"
	if *config_filename_flag != "" {

		configFileName = *config_filename_flag

	}

	configFile, err := os.Open(configFileName)
	_ = configFile
	if err != nil {
		fmt.Printf("%v\n", "Config File not found")
	} else {
		fmt.Printf("%v\n", "Config File found")
	}
}

func TestParseConfigFile() {
}

func command_config_print() {
	fmt.Printf("%v\n", "Config File Parsed As:\n\tDataPoint:\tnil\n")
}

func TestCommand_config() {
}

func main() {
	kingpin.Version("0.0.1")
	parsedFlagsAndCommands := kingpin.Parse()
	parseConfigFile()

	switch parsedFlagsAndCommands {

	case "config print":
		command_config_print()

	}
}
