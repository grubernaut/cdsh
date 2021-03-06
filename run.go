package main

import (
	"container/list"
	"fmt"

	"os"

	"github.com/grubernaut/gdsh/dsh"
	consul "github.com/hashicorp/consul/api"
	"gopkg.in/urfave/cli.v2"
)

func run(c *cli.Context) error {
	// If a service isn't requested, exit
	if c.NArg() < 1 {
		fmt.Printf("Error: Remote command not specified\n")
		return cli.ShowAppHelp(c)
	}

	// Create ExecOpts
	opts := defaultDSHConfig()
	// Set opts.Verbose output
	opts.Verbose = false
	if c.Bool("verbose") {
		fmt.Printf("Verbose flag on\n")
		opts.Verbose = true
	}

	// Find Consul server, env-var takes priority
	consulServer := c.String("server")
	if os.Getenv("CONSUL_SERVER") != "" {
		consulServer = os.Getenv("CONSUL_SERVER")
	}
	// Can't be empty, we need servers
	if consulServer == "" {
		fmt.Printf("Error: consul-server not supplied\n")
		return cli.ShowAppHelp(c)
	}

	// Create a consul catalog client
	catalog, err := consulCatalog(consulServer)
	if err != nil {
		return cli.Exit(fmt.Sprintf(
			"Error creating consul agent: %s\n", err,
		), 1)
	}

	// Parse requested service, if empty return a list of available services
	service := c.String("service")
	if service == "" {
		fmt.Printf("No service specified. Available services:\n")
		avail, err := consulServices(catalog)
		if err != nil {
			return cli.Exit(fmt.Sprintf(
				"Error querying Consul services: %s\n", err,
			), 1)
		}
		for k := range avail {
			fmt.Printf("%s\n", k)
		}
		return nil
	}

	// Add consul services to linked list
	machineList, err := populateList(catalog, service, c.String("user"))
	if err != nil {
		return cli.Exit(fmt.Sprintf(
			"Error populating DSH machine list: %s\n", err,
		), 1)
	}
	opts.MachineList = machineList

	// Set remote commands to all trailing args
	for _, v := range c.Args().Slice() {
		// Initialize remote command
		if opts.RemoteCommand == "" && v != "" {
			opts.RemoteCommand = v
			continue
		}
		opts.RemoteCommand = fmt.Sprintf("%s %s", opts.RemoteCommand, v)
	}

	if machineList.Len() < 1 {
		return cli.Exit(fmt.Sprintf(
			"No servers found for service %s", service), 1)
	}

	// Execute DSH!
	if err := opts.Execute(); err != nil {
		return cli.Exit(fmt.Sprintf("Error executing: %s", err), 1)
	}
	return nil
}

// Default GDSH config
// TODO: Make these configurable
func defaultDSHConfig() dsh.ExecOpts {
	opts := dsh.ExecOpts{
		ConcurrentShell: true,
		RemoteShell:     "ssh",
		ShowNames:       true,
	}
	return opts
}

// Returns all available consul services
func consulServices(catalog *consul.Catalog) (map[string][]string, error) {
	services, _, err := catalog.Services(nil)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// Returns a Consul Client
func consulCatalog(server string) (*consul.Catalog, error) {
	// Create Consul Client
	config := consul.DefaultConfig()
	config.Address = server
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client.Catalog(), nil
}

// Populates doubly linked machine list, with a list of requested consul services's addresses
func populateList(catalog *consul.Catalog, service string, user string) (*list.List, error) {
	services, _, err := catalog.Service(service, "", nil)
	if err != nil {
		return nil, fmt.Errorf("Error querying consul services: %s", err)
	}

	serviceList := list.New()
	for _, v := range services {
		remoteAddr := v.Address
		if user != "" {
			remoteAddr = fmt.Sprintf("%s@%s", user, remoteAddr)
		}
		addList(serviceList, remoteAddr)
	}

	return serviceList, nil
}

// Populates linked list with a supplied string element, ensuring no duplicates or nil values are stored
func addList(llist *list.List, elem string) {
	if elem == "" {
		return
	}
	if llist.Len() == 0 {
		llist.PushFront(elem)
		return
	}
	// Verify no items match currently
	for e := llist.Front(); e != nil; e = e.Next() {
		if e.Value == elem {
			return
		}
	}
	llist.PushBack(elem)
}
