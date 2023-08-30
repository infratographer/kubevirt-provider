package server

import (
	"errors"
	"net/http"

	// "slices"
	"strings"

	"go.infratographer.com/x/events"
	"go.infratographer.com/x/gidx"
	"golang.org/x/exp/slices"

	"go.infratographer.com/kubevirt-provider/internal/virtualmachine"
)

func (s *Server) ProcessChange(messages <-chan events.Message[events.ChangeMessage]) {
	var vm *virtualmachine.VirtualMachine

	for msg := range messages {
		m := msg.Message()

		// m.
		if slices.ContainsFunc(m.AdditionalSubjectIDs, s.LocationCheck) || len(s.Locations) == 0 {
			if m.EventType != string(events.DeleteChangeType) {
				_, err := virtualmachine.NewVirtualMachine(s.Context, s.Logger, s.APIClient, m.SubjectID, m.AdditionalSubjectIDs)
				if err != nil {
					// s.Logger.Errorw("unable to initialize virtualmachine", "error", err, "messageID", msg.UUID, "message", msg.Payload)

					if errors.Is(err, http.ErrNoLocation) {
						// ack and ignore
						msg.Ack()
					} else {
						// nack and retry
						// msg.Nak()
						return // XXX
					}

					continue
				}
			} else {
				vm = &virtualmachine.VirtualMachine{
					VirtualMachineID: m.SubjectID,
					VMType:           virtualmachine.TypeVM,
				}
			}

			if vm != nil && vm.VMType != virtualmachine.TypeNoVM {
				switch {
				case m.EventType == string(events.CreateChangeType) && vm.VMType == virtualmachine.TypeVM:
					s.Logger.Debugw("requesting address for virtualmachine", "virtualmachine", vm.VirtualMachineID.String())

					if err := s.processVirtualMachineChangeCreate(vm); err != nil {
						s.Logger.Errorw("handler unable to request address for virtualmachine", "error", err, "virtualmachine", vm.VirtualMachineID.String())
						// msg.Nack()
					}
				case m.EventType == string(events.DeleteChangeType) && vm.VMType == virtualmachine.TypeVM:
					s.Logger.Debugw("releasing address from virtualmachine", "virtualmachine", vm.VirtualMachineID.String())

					if err := s.processVirtualMachineChangeDelete(vm); err != nil {
						s.Logger.Errorw("handler unable to release address from virtualmachine", "error", err, "virtualmachine", vm.VirtualMachineID.String())
						// msg.Nack()
					}
				default:
					// s.Logger.Debugw("Ignoring event", "virtualmachine", vm.VirtualMachineID.String(), "message", msg.Payload)
				}
			}
		}
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

func (s *Server) processVirtualMachineChangeCreate(vm *virtualmachine.VirtualMachine) error {
	// // for now, limit to one IP address per virtualmachine
	// if len(vm.VMData.VirtualMachine.IPAddresses) == 0 {
	// 	if ip, err := ipam.RequestAddress(s.Context, s.IPAMClient, s.Logger, s.IPBlock, vm.VirtualMachineID.String(), vm.VmData.VirtualMachine.Owner.ID); err != nil {
	// 		return err
	// 	} else {
	// 		msg := events.EventMessage{
	// 			EventType: "ip-address.assigned",
	// 			SubjectID: vm.VirtualMachineID,
	// 			Timestamp: time.Now().UTC(),
	// 		}

	// 		if err := s.Publisher.PublishEvent(s.Context, "load-balancer", msg); err != nil {
	// 			s.Logger.Debugw("failed to publish event", "error", err, "ip", ip, "virtualmachine", vm.VirtualMachineID, "block", s.IPBlock)
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func (s *Server) processVirtualMachineChangeDelete(vm *virtualmachine.VirtualMachine) error {
	// if err := ipam.ReleaseAddress(s.Context, s.IPAMClient, s.Logger, vm.VirtualMachineID.String()); err != nil {
	// 	return err
	// }

	// msg := events.EventMessage{
	// 	EventType: "ip-address.unassigned",
	// 	SubjectID: vm.VirtualMachineID,
	// 	Timestamp: time.Now().UTC(),
	// }

	// if err := s.Publisher.PublishEvent(s.Context, "load-balancer", msg); err != nil {
	// 	s.Logger.Debugw("failed to publish event", "error", err, "virtualmachine", vm.VirtualMachineID, "block", s.IPBlock)
	// 	return err
	// }

	return nil
}

func (s *Server) LocationCheck(i gidx.PrefixedID) bool {
	for _, s := range s.Locations {
		if strings.HasSuffix(i.String(), s) {
			return true
		}
	}

	return false
}
