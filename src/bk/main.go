package main

import (
	"beastkeeper/src/bk/instanceTypes"
	"beastkeeper/src/bk/states"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
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
	installImageManager            InstallImageManager
)

// InstanceStateMachine has an "Enforce" function that iterates over an Instance
// struct and execute local or remote (SSH) commands that should move the state
// machine closer to the desired state. The state machine is defined as
// collection of State structs, containing two functions; one to test the state
// and a matching function that should cause the state to be true after it's
// successful execution.
type InstanceStateMachine struct {
	states   []states.T_State
	instance instanceTypes.BaseInstance
}

func (self *InstanceStateMachine) GenerateStates() {
	if self.instance.Type == instanceTypes.VM {
		diskImageExists := &states.DiskImageExistsState{BaseState: states.BaseState{}}
		diskImageExists.SetMaxAttempts(5)
		self.states = append(self.states, diskImageExists)
	}
}

func (self *InstanceStateMachine) Enforce() bool {

	self.GenerateStates()

	fmt.Println("states generated", len(self.states))

	for _, state := range self.states {
		fmt.Println("assesing")
		state.SetAttempts(0)
		for !state.Assess(self.instance) && (state.GetAttempts() < state.GetMaxAttempts()) {
			fmt.Printf("attempt:%d of %d\n", state.GetAttempts(), state.GetMaxAttempts())
			state.Enforce(self.instance)
			state.Advance()
		}
	}
	return true
}

// BeastKeeperConfiguration structs describe a single full configuration of the
// BeastKeeper application. There will initially only be one instance of this
// type held in a global variable; but others may be introduced later if we add
// features to manage multiple configurations
type BeastKeeperConfiguration struct {
	Instances []instanceTypes.BaseInstance
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

type InstallImage struct {
	version      FloatType
	architecture int
}

type InstallImageManager struct {
	InstallImagePath   string
	ImageDownloadQueue []InstallImage
}

func (self *InstallImageManager) ImageExists() bool {
	return false
}

func (self *InstallImageManager) ImageHasCorrectCheckSum() bool {
	return false
}

func (self *InstallImageManager) RegisterImageForDownload() {
}

func (self *InstallImageManager) ProcessImageDownloadQueue() {
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

func enforceInstanceConfig(instance instanceTypes.BaseInstance, channel chan string) {

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
