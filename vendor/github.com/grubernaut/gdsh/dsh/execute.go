package dsh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Execute remote commands for each host
func (e *ExecOpts) Execute() error {
	signals := make(chan Signal)
	for elem := e.MachineList.Front(); elem != nil; elem = elem.Next() {
		cmdOpts := e.buildCmdOpts(elem.Value.(string))
		if e.Verbose {
			fmt.Printf("Dumping parameters passed to exec\n")
			fmt.Printf("%#v\n", cmdOpts)
		}
		// Spawn a new goroutine
		go executeShell(e.RemoteShell, cmdOpts, signals, e.ShowNames, elem.Value.(string))
	}

	// Block until routines are cleaned up
	var err error
	for i := 0; i < e.MachineList.Len(); i++ {
		select {
		case signal := <-signals:
			if signal.err != nil {
				fmt.Printf("Error executing: %s\n", signal.errOutput)
				err = signal.err
			}
		}
	}
	return err
}

// Build up command options for each item in linked list
func (e *ExecOpts) buildCmdOpts(elem string) []string {
	var opts []string
	if e.RemoteCommandOpts != "" {
		opts = append(opts, e.RemoteCommandOpts)
	}
	// split machine name and username
	if strings.Contains(elem, "@") {
		userHost := strings.Split(elem, "@")
		opts = append(opts, "-l")
		opts = append(opts, userHost...)
	} else {
		opts = append(opts, elem)
	}
	opts = append(opts, e.RemoteCommand)
	return opts
}

// Performs the actual execution of the Remote Shell command
func executeShell(cmd string, cmdOpts []string, c chan Signal, names bool, name string) {
	// hopefully you don't need it
	var errOutput bytes.Buffer
	run := exec.Command(cmd, cmdOpts...)
	run.Stderr = io.Writer(&errOutput)
	stdout, err := run.StdoutPipe()
	if err != nil {
		c <- Signal{
			err: err,
		}
		return
	}
	run.Env = os.Environ()

	// Create a new scanner from stdout pipe
	scanner := bufio.NewScanner(stdout)
	// While we have stdout to print, print it.
	go func() {
		for scanner.Scan() {
			if names {
				fmt.Printf("%s: %s\n", name, scanner.Text())
			} else {
				fmt.Printf("%s\n", scanner.Text())
			}
		}
	}()

	if err := run.Start(); err != nil {
		c <- Signal{
			err:       err,
			errOutput: errOutput.String(),
		}
		return
	}

	// Block for command to finish
	if err := run.Wait(); err != nil {
		c <- Signal{
			err:       err,
			errOutput: errOutput.String(),
		}
		return
	}

	// Non-error case
	c <- Signal{nil, ""}
}
