package virtualmachine

import "go.infratographer.com/x/gidx"

const (
	VMPrefix = "virtmachine"

	// TypeLB indicates that the subject of a message is a virtualmachine
	TypeVM = 1
	// TypeAssocLB indicates that the virtualmachine was found in associated subjects
	TypeAssocVM = 2
	// TypeNoLB indicates that a virtualmachine was not found in the message
	TypeNoVM = 0
)

type VirtualMachine struct {
	VirtualMachineID gidx.PrefixedID
	// VMData         *lbapi.GetVirtualMachine
	VMType int
}
