package virtualmachine

import "go.infratographer.com/x/gidx"

const (
	// VMPrefix indicates the prefix
	VMPrefix = "virtmachine"

	// TypeVM indicates that the subject of a message is a virtualmachine
	TypeVM = 1
	// TypeAssocVM indicates that the virtualmachine was found in associated subjects
	TypeAssocVM = 2
	// TypeNoVM indicates that a virtualmachine was not found in the message
	TypeNoVM = 0
)

type VirtualMachine struct {
	VirtualMachineID gidx.PrefixedID
	// VMData         *lbapi.GetVirtualMachine
	VMType int
}
