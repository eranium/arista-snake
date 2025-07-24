package main

import (
	"errors"
	"github.com/aristanetworks/goeapi"
)

type InterfaceLEDStatus int

const (
	InterfaceLEDStatusOff InterfaceLEDStatus = iota
	InterfaceLEDStatusGreen
	InterfaceLEDStatusOrange
)

type InterfaceLEDController interface {
	SetStatus(intf string, s InterfaceLEDStatus) error
	SetAllStatuses(s InterfaceLEDStatus) error
}

type EnableCommand struct{}

func (c *EnableCommand) GetCmd() string {
	return "enable"
}

type ConfigureCommand struct{}

func (c *ConfigureCommand) GetCmd() string {
	return "configure"
}

type InterfaceCommand struct {
	intf string
}

func (c *InterfaceCommand) GetCmd() string {
	return "interface " + c.intf
}

type ShutdownCommand struct {
	shut bool
}

func (c *ShutdownCommand) GetCmd() string {
	cmd := "shutdown"

	if c.shut {
		cmd = "no " + cmd
	}

	return cmd
}

type TrafficLoopBackCommand struct {
	enable bool
}

func (c *TrafficLoopBackCommand) GetCmd() string {
	cmd := "traffic-loopback source system device mac"

	if !c.enable {
		cmd = "no " + cmd
	}

	return cmd
}

type AristaInterfaceLEDController struct {
	managedIntfs []string
	node         *goeapi.Node
}

func NewAristaInterfaceLEDController(managedIntfs []string, node *goeapi.Node) *AristaInterfaceLEDController {
	return &AristaInterfaceLEDController{
		managedIntfs: managedIntfs,
		node:         node,
	}
}

func (a *AristaInterfaceLEDController) SetStatus(intf string, s InterfaceLEDStatus) error {
	handle, err := a.node.GetHandle("json")
	if err != nil {
		return err
	}

	_ = handle.AddCommand(&EnableCommand{})
	_ = handle.AddCommand(&ConfigureCommand{})
	_ = handle.AddCommand(&InterfaceCommand{intf: intf})

	switch s {
	case InterfaceLEDStatusOff:
		_ = handle.AddCommand(&ShutdownCommand{shut: false})
		_ = handle.AddCommand(&TrafficLoopBackCommand{enable: false})
	case InterfaceLEDStatusOrange:
		_ = handle.AddCommand(&ShutdownCommand{shut: true})
		_ = handle.AddCommand(&TrafficLoopBackCommand{enable: false})
	case InterfaceLEDStatusGreen:
		_ = handle.AddCommand(&ShutdownCommand{shut: false})
		_ = handle.AddCommand(&TrafficLoopBackCommand{enable: true})
	default:
		return errors.New("unknown interface LED status")
	}

	return handle.Call()
}

func (a *AristaInterfaceLEDController) SetAllStatuses(s InterfaceLEDStatus) error {
	for _, intf := range a.managedIntfs {
		if err := a.SetStatus(intf, s); err != nil {
			return err
		}
	}

	return nil
}

type NoOpInterfaceLEDController struct{}

func NewNoOpInterfaceLEDController() NoOpInterfaceLEDController {
	return NoOpInterfaceLEDController{}
}

func (n NoOpInterfaceLEDController) SetStatus(_ string, _ InterfaceLEDStatus) error {
	return nil
}

func (n NoOpInterfaceLEDController) SetAllStatuses(_ InterfaceLEDStatus) error {
	return nil
}
