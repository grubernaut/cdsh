package dsh

import (
	"bufio"
	"container/list"
	"os"
)

// Generic functions for building the linked machine list

// ReadMachineList will read a machine list from file, and append it
// to the current machine list
func (e *ExecOpts) ReadMachineList(path string) error {
	// Read file path
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scan := bufio.NewScanner(file)
	// Add read-in file to machineList
	for scan.Scan() {
		e.MachineList = addList(e.MachineList, scan.Text())
	}

	return nil
}

// SplitListAndAdd will split a comma-deliminated list of machines
// and add them to the linked list
func (e *ExecOpts) SplitListAndAdd(input string) error {
	return nil
}

func addList(llist *list.List, elem string) *list.List {
	// if empty list, push elem to front
	if llist.Len() == 0 {
		llist.PushFront(elem)
		return llist
	}
	// make sure element isn't nil
	if elem == "" {
		return llist
	}
	// loop through list; see if there's a match currently; if so return early
	for e := llist.Front(); e != nil; e = e.Next() {
		if e.Value == elem {
			return llist
		}
	}
	// Element not in linked list, add to end of list
	if elem != "" {
		llist.PushBack(elem)
	}
	// No error, hopefully
	return llist
}
