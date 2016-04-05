package instanceTypes

import (
	"encoding/json"
	"errors"
	"github.com/twinj/uuid"
	"net"
	"path/filepath"
	"strings"
)

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

// Type and Enum construct for describing Instance types
type InstanceType int

const (
	VM InstanceType = iota
	BM InstanceType = iota
)

// Overriding the MarshalJSON method of our InstanceType so we can use our enum
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

// Overriding the UnmarshalJSON method of our InstanceType so we can use our enum
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

// BaseInstance structs contain the data required to describe an individual FreeBSD
// instance deployed anywhere. This can be bare-metal, or virtual machine either
// local or at a provider.  Application containers such as jetpack pods are not
// included in this, and have their own data structure
type BaseInstance struct {
	ID                      *UUID `json:",UUID"`
	Label                   string
	Type                    InstanceType
	OSVersion               float32
	IsRunning               bool
	Path                    string
	RootDiskImageSize       string
	Address                 net.IP
	AdminLogin              string
	SSHPort                 int
	SSHKeyPassphraseEnabled bool
	Containers              []BaseApplicationContainerInstance
}

// Constructor method for BaseInstance structs
func NewBaseInstance(args ...interface{}) *BaseInstance {
	newInstance := &BaseInstance{ID: new(UUID)}

	for argNum, arg := range args {

		switch argNum {
		case 0:
			uuid, _ := uuid.Parse(arg.(string))
			newInstance.ID.UUID = uuid
		case 1:
			newInstance.Label = arg.(string)
		case 2:
			newInstance.Type = arg.(InstanceType)
		case 3:
			newInstance.OSVersion = arg.(float32)
		case 4:
			newInstance.IsRunning = arg.(bool)
		case 5:
			newInstance.Path = arg.(string)
		case 6:
			newInstance.RootDiskImageSize = arg.(string)
		case 7:
			newInstance.Address = arg.(net.IP)
		case 8:
			newInstance.AdminLogin = arg.(string)
		case 9:
			newInstance.SSHPort = arg.(int)
		case 10:
			newInstance.SSHKeyPassphraseEnabled = arg.(bool)
		case 11:
			newInstance.Containers = arg.([]BaseApplicationContainerInstance)
		}
	}

	return newInstance

}

func (self *BaseInstance) GetDiskImageFileName() string {
	return filepath.Join(self.Path, self.Label+".img")
}
