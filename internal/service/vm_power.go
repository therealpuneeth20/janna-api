package service

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	vmware_types "github.com/vmware/govmomi/vim25/types"

	"github.com/vterdunov/janna-api/internal/types"
)

const (
	off       = vmware_types.VirtualMachinePowerStatePoweredOff
	on        = vmware_types.VirtualMachinePowerStatePoweredOn
	suspended = vmware_types.VirtualMachinePowerStateSuspended
)

// Power changes VM power state
func (s *service) VMPower(ctx context.Context, params *types.VMPowerParams) error {
	vm, err := findByUUID(ctx, s.Client, params.Datacenter, params.UUID)
	if err != nil {
		return err
	}

	switch params.State {
	case "on":
		err = powerOn(ctx, vm)
	case "off":
		err = powerOff(ctx, vm)
	case "suspend":
		err = suspend(ctx, vm)
	case "reboot":
		err = reboot(ctx, vm)
	case "reset":
		err = reset(ctx, vm)
	}

	return err
}

// powerOn power on Virtual Machine
func powerOn(ctx context.Context, vm *object.VirtualMachine) error {
	state, err := getVMPowerState(ctx, vm)
	if err != nil {
		return err
	}

	switch state {
	case on:
		return nil

	case off, suspended:
		task, err := vm.PowerOn(ctx)
		if err != nil {
			return err
		}

		return task.Wait(ctx)

	default:
		// actually there are no more states. safe to the future.
		return fmt.Errorf("could not power on Virtual Machine from %s state", state)
	}
}

// powerOff power off Virtual Machine
func powerOff(ctx context.Context, vm *object.VirtualMachine) error {
	state, err := getVMPowerState(ctx, vm)
	if err != nil {
		return err
	}

	switch state {
	case off:
		return nil

	case on:
		err := vm.ShutdownGuest(ctx)
		if err != nil && isToolsUnavailable(err) {
			task, powerErr := vm.PowerOff(ctx)
			if powerErr != nil {
				return powerErr
			}

			return task.Wait(ctx)
		}
		return err

	case suspended:
		task, err := vm.PowerOff(ctx)
		if err != nil {
			return err
		}

		return task.Wait(ctx)

	default:
		return fmt.Errorf("could not power off Virtual Machine from %s state", state)
	}
}

// reboot Virtual Machine. It tries to use VMWareTools to call guest agent to reboot the VM.
// And as the last way, the method tries to reset VM.
func reboot(ctx context.Context, vm *object.VirtualMachine) error {
	state, err := getVMPowerState(ctx, vm)
	if err != nil {
		return err
	}

	if state != on {
		return fmt.Errorf("could not reboot Virtual Machine from %s state", state)
	}

	err = vm.RebootGuest(ctx)
	if err != nil && isToolsUnavailable(err) {
		task, resetErr := vm.Reset(ctx)
		if resetErr != nil {
			return resetErr
		}

		return task.Wait(ctx)
	}

	return err
}

// reset Virtual Machine
func reset(ctx context.Context, vm *object.VirtualMachine) error {
	state, err := getVMPowerState(ctx, vm)
	if err != nil {
		return err
	}

	if state != on {
		return fmt.Errorf("could not reset Virtual Machine from %s state", state)
	}

	task, err := vm.Reset(ctx)
	if err != nil {
		return err
	}

	return task.Wait(ctx)
}

// suspend Virtual Machine
func suspend(ctx context.Context, vm *object.VirtualMachine) error {
	state, err := getVMPowerState(ctx, vm)
	if err != nil {
		return err
	}

	switch state {
	case suspended:
		return nil

	case on:
		task, err := vm.Suspend(ctx)
		if err != nil {
			return err
		}

		return task.Wait(ctx)

	default:
		return fmt.Errorf("could not suspend Virtual Machine from %s state", state)
	}
}

func getVMPowerState(ctx context.Context, vm *object.VirtualMachine) (vmware_types.VirtualMachinePowerState, error) {
	state, err := vm.PowerState(ctx)
	if err != nil {
		return "", err
	}

	return state, err
}

func isToolsUnavailable(err error) bool {
	if soap.IsSoapFault(err) {
		soapFault := soap.ToSoapFault(err)
		if _, ok := soapFault.VimFault().(vmware_types.ToolsUnavailable); ok {
			return ok
		}
	}

	return false
}
