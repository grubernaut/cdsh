package dsh

import "container/list"

//ExecOpts does things
type ExecOpts struct {
	MachineList       *list.List
	CommandList       *list.List
	ShowNames         bool
	RemoteShell       string
	RemoteCommand     string
	RemoteCommandOpts string
	ConcurrentShell   bool
	Verbose           bool
}

// Signal is returned from a goroutine via a channel
type Signal struct {
	err       error
	errOutput string
}
