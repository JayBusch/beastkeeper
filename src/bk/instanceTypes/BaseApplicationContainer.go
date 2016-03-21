package instanceTypes

import (
	"github.com/twinj/uuid"
	"net"
)

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
