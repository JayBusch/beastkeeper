package main

import (
	"fmt"
	"github.com/twinj/uuid"
	"gopkg.in/alecthomas/kingpin.v2"
	"net"
	"os"
)

// This var block contains the command line commands and flags for BeastKeeper.
var (
	command_config_flag       = kingpin.Command("config", "work with the config file")
	command_config_print_flag = command_config_flag.Command("print", "print the config file")
	config_filename_flag      = kingpin.Flag("configFileName", "set the config file name").Short('c').String()
)

// This var block contains globally declared variables for BeastKeeper.
var (
	beastKeeperMasterConfiguration BeastKeeperConfiguration
)

//Type and Enum construct for describing Instance types
type InstanceType int

const (
	VM  InstanceType = iota
	POD InstanceType = iota
)

// Instance structs contain the data required to describe an individual FreeBSD instance deployed anywhere
// This can be bare-metal, or virtual machine either local or at a provider.
// Application containers such as jetpack pods are not included in this, and have their own data structure
type Instance struct {
	ID         uuid.UUID
	Label      string
	Type       InstanceType
	address    net.IP
	containers []ApplicationContainerInstance
}

//Type and Enum construct for describing ApplicationContainer types
type ApplicationContainerType int

const (
	JetPack ApplicationContainerType = iota
	Docker  ApplicationContainerType = iota
)

// ApplicationContainer structs contain the data requied to describe isntances of OS level virtualized
// application containers such as jetpack pods.
type ApplicationContainerInstance struct {
	ID      uuid.UUID
	Label   string
	Type    ApplicationContainerType
	address net.IP
}

// BeastKeeperConfiguration structs describe a single full configuration of the BeastKeeper application
// There will initially only be one instance of this type held in a global variable; but others may
// be introduced later if we add features to manage multiple configurations
type BeastKeeperConfiguration struct {
	instances []InstanceType
}

// parseConfigFile reads a config file and marshalls the JSON data to a GO struct representing
// the configuration of BeastKeeper.
// The User can speficy the config file name on the command line with the -c or --configFileName flags.
// If no file is specified, the defailt value of "config.bk" is used.
// BeastKeeper looks in the current directory for that default file.
func parseConfigFile(configFileName string) {

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

	configFileName := "config.bk"
	if *config_filename_flag != "" {

		configFileName = *config_filename_flag

	}

	parseConfigFile(configFileName)

	switch parsedFlagsAndCommands {

	case "config print":
		command_config_print()

	}
}
