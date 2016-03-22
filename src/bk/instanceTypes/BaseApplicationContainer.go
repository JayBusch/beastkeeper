package instanceTypes

import (
	"github.com/twinj/uuid"
	"net"
)

//Type and Enum construct for describing ApplicationContainer types
type BaseApplicationContainerType int

const (
	JetPack BaseApplicationContainerType = iota
	Docker  BaseApplicationContainerType = iota
)

// ApplicationContainer structs contain the data requied to describe isntances
// of OS level virtualized application containers such as jetpack pods.
type BaseApplicationContainerInstance struct {
	ID      uuid.UUID
	Label   string
	Type    BaseApplicationContainerType
	Address net.IP
}
