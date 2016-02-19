package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/twinj/uuid"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
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

// Type and Enum construct for describing Instance types
type InstanceType int

const (
	VM InstanceType = iota
	BM InstanceType = iota
)

// Overriding the UnmarshalJSON method of our InstaceType so we can use our enum
func (self InstanceType) UnmarshalJSON(b []byte) error {
	switch strings.Trim(string(b), "\"") {
	case "VM":
		self = VM
		return nil
	case "BM":
		self = BM
		return nil
	default:
		return errors.New("Un-Recognized InstanceType")
	}
}

// The UUID struct is created to hold a single UUID type such that it's
// UnmarshalJSON method can be overriden in order to parse the UUID during JSON
// marshalling
type UUID struct {
	UUID uuid.UUID
}

func (self *UUID) UnmarshalJSON(b []byte) error {

	s := strings.Trim(string(b), "\"")
	uuid, uuidErr := uuid.Parse(s)
	self.UUID = uuid
	if self.UUID == nil || uuidErr != nil {
		return errors.New("Could not parse UUID")
	}
	return nil
}

// Instance structs contain the data required to describe an individual FreeBSD
// instance deployed anywhere. This can be bare-metal, or virtual machine either
// local or at a provider.  Application containers such as jetpack pods are not
// included in this, and have their own data structure
type Instance struct {
	ID         *UUID `json:",UUID"`
	Label      string
	Type       InstanceType
	Address    net.IP
	AdminLogin string
	Containers []ApplicationContainerInstance
}

//Type and Enum construct for describing ApplicationContainer types
type ApplicationContainerType int

const (
	JetPack ApplicationContainerType = iota
	Docker  ApplicationContainerType = iota
)

// ApplicationContainer structs contain the data requied to describe isntances
// of OS level virtualized application containers such as jetpack pods.
type ApplicationContainerInstance struct {
	ID      uuid.UUID
	Label   string
	Type    ApplicationContainerType
	Address net.IP
}

// BeastKeeperConfiguration structs describe a single full configuration of the
// BeastKeeper application. There will initially only be one instance of this
// type held in a global variable; but others may be introduced later if we add
// features to manage multiple configurations
type BeastKeeperConfiguration struct {
	Instances []Instance
}

// parseConfigFile reads a config file and marshalls the JSON data to a GO
// struct representing the configuration of BeastKeeper. The User can specify
// the config file name on the command line with the -c or --configFileName
// flags. If no file is specified, the defailt value of "config.bk" is used.
// BeastKeeper looks in the current directory for that default file.
func parseConfigFile(configFileName string) BeastKeeperConfiguration {

	configFileData, fileErr := ioutil.ReadFile(configFileName)

	config := new(BeastKeeperConfiguration)
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	jsonErr := json.Unmarshal(configFileData, &config)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return *config
}

// commandConfigPrint allows the user to print the currently configured config
// file to stdout.  The format is a valid JSON config file in and of itself.
// This can be used along with other command line flags to construct a permanent
// config file.
func commandConfigPrint() {

	configBytes, err := json.Marshal(beastKeeperMasterConfiguration)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(configBytes)

}

func main() {
	kingpin.Version("0.0.1")
	parsedFlagsAndCommands := kingpin.Parse()

	configFileName := "config.bk"
	if *config_filename_flag != "" {

		configFileName = *config_filename_flag

	}

	beastKeeperMasterConfiguration = parseConfigFile(configFileName)

	switch parsedFlagsAndCommands {

	case "config print":
		commandConfigPrint()

	}
}
