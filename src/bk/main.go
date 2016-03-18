package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/twinj/uuid"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	//"math/rand"
	"net"
	"os"
	"strings"
	//"time"
)

// This var block contains the command line commands and flags for BeastKeeper.
var (
	command_config_flag       = kingpin.Command("config", "work with the config file")
	command_config_print_flag = command_config_flag.Command("print", "print the config file")
	config_filename_flag      = kingpin.Flag("configFileName", "set the config file name").Short('c').String()

	command_enforce_flag = kingpin.Command("enforce", "attempt to create a system state that matches the configuration")
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

// Overriding the MarshalJSON method of our InstaceType so we can use our enum
func (self InstanceType) MarshalJSON() ([]byte, error) {
	switch self {
	case VM:
		return json.Marshal("VM")
	case BM:
		return json.Marshal("BM")
	default:
		return nil, errors.New("Un-Recognized InstanceType")
	}
}

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

// Overriding the MarshalJSON method of the UUID type so we can return a string
func (self *UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.UUID.String())
}

// Overriding the UnmarshalJSON method of the UUID type so we can parse the UUID
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
	Path       string
	Address    net.IP
	AdminLogin string
	Containers []ApplicationContainerInstance
}

// InstanceStateMachine has an "Enforce" function that iterates over an Instance
// struct and execute local or remote (SSH) commands that should move the state
// machine closer to the desired state. The state machine is defined as
// collection of State structs, containing two functions; one to test the state
// and a matching function that should cause the state to be true after it's
// successful execution.
type InstanceStateMachine struct {
	states   []t_State
	instance Instance
}

func (self *InstanceStateMachine) GenerateStates() {

	//self.states := []interface{}

	if self.instance.Type == VM {
		fmt.Println("len(self.states): ", len(self.states))
		diskImageExists := &DiskImageExistsState{State: State{maxAttempts: 5}}
		self.states = append(self.states, diskImageExists)
		fmt.Println("len(self.states): ", len(self.states))
	}

}

type t_State interface {
	assess() bool
	advance()
	getAttempts() int
	setAttempts(int)
	getMaxAttempts() int
	setMaxAttempts(int)
}

type State struct {
	attempts    int
	maxAttempts int
}

func (self *State) advance() {
	self.attempts = self.attempts + 1
}
func (self State) assess() bool {
	return true
}

func (self State) getAttempts() int {
	return self.attempts
}

func (self *State) setAttempts(attempts int) {
	self.attempts = attempts
}

func (self State) getMaxAttempts() int {
	return self.maxAttempts
}

func (self *State) setMaxAttempts(maxAttempts int) {
	self.maxAttempts = maxAttempts
}

type DiskImageExistsState struct {
	State
}

//func (self DiskImageExistsState) attempts() int {
//	return 0
//}

func (self DiskImageExistsState) assess() bool {
	return false
}

func (self *InstanceStateMachine) Enforce() bool {

	self.GenerateStates()

	fmt.Println("states generated", len(self.states))

	for _, state := range self.states {
		fmt.Println("assesing")
		state.setAttempts(0)
		for !state.assess() && (state.getAttempts() < state.getMaxAttempts()) {
			fmt.Printf("attempt:%d of %d\n", state.getAttempts(), state.getMaxAttempts())
			state.advance()
		}
	}
	return true
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
	beastKeeperMasterConfiguration = *config

	//return here in order to increase code testability
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

// commandEnforce causes BeastKeeper to attempt creating a state that matches
// it's configuration.
func commandEnforce() {

	// First we need to spawn a goroutine for each instance.
	// Those respective goroutines will be responsible for assesing state
	// of each Instance, and performing actions to achieve the desired
	// state.  All of these goroutines should post status messages back to
	// this main thread which will output them to the user.  We also use
	// these channel messages to determine state; a message of
	// "config enforced" indicates that the goroutine is finished and will
	// exit.  Once all goroutines have exited this function will return.

	var instanceChannels = []chan string{}

	for _, instance := range beastKeeperMasterConfiguration.Instances {
		instanceChannel := make(chan string)
		instanceChannels = append(instanceChannels, instanceChannel)
		go enforceInstanceConfig(instance, instanceChannel)
	}

	completedRoutines := 0

	for completedRoutines < len(instanceChannels) {
		for _, channel := range instanceChannels {
			message := <-channel
			if message == "config enforced" {
				completedRoutines++
				fmt.Printf("%v completedRoutines\n", completedRoutines)
			}
		}
	}

}

func enforceInstanceConfig(instance Instance, channel chan string) {
	//time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)

	// For this Instance we need to determine it's state.

	ism := InstanceStateMachine{instance: instance}

	ism.Enforce()

	channel <- "config enforced"
}

// General order of operations here is:
//
// 1. Parse our command line arguments and set config variables accordingly
// 2. Parse our configuration file
// 3. Iterate through the instances as defined in config, and launch an Instance
//    State Machine for each one.
// 4. When all Instance State Machines have reached either True or Error states
//    report to STDOUT and exit.
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
		commandConfigPrint()
	case "enforce":
		commandEnforce()
	}

}
