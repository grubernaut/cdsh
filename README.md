Consul DSH
==========

`go get -u github.com/grubernaut/cdsh`

The remote shell `exec` for consul services.

### Quick How-to

```
cdsh --server <consul-server-address> --service <consul-service> --user <username> '<bash command to execute>'
```

**Note:** 
Be sure that you're actually able to SSH to the target consul service. (SSH Keys, Host key verification, etc)

### Options

* `--server` / `-s`: Consul server to query against (Required)
* `--service` / `-S`: Target service to run command against
* `--user` / `-u`: Remote user to connect as. If unspecified, defaults to current local user

More options are likely to come in the future, such as `--tag` `--token`, `--node`, etc...

If `--service` is left unspecified, `cdsh` will return a list of all available services. This will likely be changed
in the future to a separate flag entirely.

### Doesn't Consul have a native `exec` command?

While consul _does_ have a native `exec` command, there are a couple issues with using `consul exec`.

> Agents are informed about the new job using the event system, which propagates messages via the gossip protocol. As a 
> result, delivery is best-effort, and there is no guarantee of execution.
> ...
> While events are purely gossip driven, remote execution relies on the KV store as a message broker.  

Using Consul K/V as a message broker _definitely_ has it's benefits, but there are often times in daily operations work,
where a user might want to know immediately if they cannot connect to a service.

However, the main concern with using `consul exec` is as follows:

> Verbose output warning: use care to make sure that your command does not produce a large volume of output. Writes to
> the KV store for this output go through the Consul servers and the Raft consensus algorithm, so having a large number
> of nodes in the cluster flow a large amount of data through the KV store could make the cluster unavailable.

There are quite a few times where I've needed to fetch fairly verbose log output from N-servers, and pipe the output to
a file for further analysis. Using DSH completely abstracts the command output from Consul K/V, allowing the user to run
fairly verbose commands without worry.
