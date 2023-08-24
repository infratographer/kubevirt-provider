// Package virtualmachine provides functions and types for inspecting virtualmachines
package virtualmachine

import (
	"context"
	"net/http"

	"go.infratographer.com/x/gidx"
	"go.uber.org/zap"
)

// NewVirtualMachine will create a new virtualmachine object
func NewVirtualMachine(ctx context.Context, logger *zap.SugaredLogger, client *http.Client, subj gidx.PrefixedID, adds []gidx.PrefixedID) (*VirtualMachine, error) {
	v := new(VirtualMachine)
	v.isVirtualMachine(subj, adds)

	// if v.VMType != TypeNoVM {
	// 	_, err := client.GetVirtualMachine(ctx, v.VirtualMachineID.String())
	// 	if err != nil {
	// 		logger.Errorw("unable to get virtualmachine from API", "error", err)

	// 		return nil, err
	// 	}

	// 	// v. = data
	// }

	// Stubbed out for now
	v.VirtualMachineID = gidx.MustNewID(VMPrefix)
	v.VMType = TypeVM
	return v, nil
}

func (l *VirtualMachine) isVirtualMachine(subj gidx.PrefixedID, adds []gidx.PrefixedID) {
	check, subs := getVMFromAddSubjs(adds)

	switch {
	case subj.Prefix() == VMPrefix:
		l.VirtualMachineID = subj
		l.VMType = TypeVM

		return
	case check:
		l.VirtualMachineID = subs
		l.VMType = TypeAssocVM

		return
	default:
		l.VMType = TypeNoVM
		return
	}
}

func getVMFromAddSubjs(adds []gidx.PrefixedID) (bool, gidx.PrefixedID) {
	for _, i := range adds {
		if i.Prefix() == VMPrefix {
			return true, i
		}
	}

	id := new(gidx.PrefixedID)

	return false, *id
}
