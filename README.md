# Usage - nflow-generator

This program generates mock netflow (v5) data that can be used to test netflow collector programs. 
The program simulates a router that is exporting flow records to the collector.
It is useful for determining whether the netflow collector is operating and/or receiving netflow datagrams.

nflow-generator generates several netflow datagrams per second, each with 8 or 16 records for varying kinds of traffic (HTTP, SSH, SNMP, DNS, MySQL, and many others.)

### Docker Image Run (Easiest)

Simply run in a container and pass any arguments at runtime. Below is an example passing the `--help` flag:

```
docker run -it --rm networkstatic/nflow-generator --help
```

To generate mock flow data simply add the target IP and port:

```
docker run -it --rm networkstatic/nflow-generator -t <ip> -p <port>
```

### Build

Install [Go](http://golang.org/doc/install), then:

	git clone https://github.com/nerdalert/nflow-generator.git -or -
	git clone https://github.com/richb-hanover/nflow-generator.git
	cd <dir>
	go build

Go build will leave a binary in the root directory that can be run.
	
### RUN

Feed it the target collector and port, and optional "false-index" flag:

	./nflow-generator -t <ip> -p <port> [ -f | --false-index ]

### Update - May 2017

The original mock netflow generator placed random values in several fields which confused 
certain netflow collectors that complained about inaccurate time stamps, 
and were confused by the random values sent in the input and output interface fields. This update:

* Sets the `SysUptime`, `unix_secs`, and `unix_nsecs` fields of the Netflow datagrams to sensible (UTC) values
* Generates a unique `flow_sequence` value for each netflow datagram
* Creates reasonable start/stop times for flows, so the First is set to (now-X) and Last to (now-Y), where X & Y are random times, and X > Y.
* If the --false-index (-f) flag is set on the command line, 
use this algorithm to set the interface indexes to 1 or 2:
If the source address > dest address, input interface is set to 1, and set to 2 otherwise,
and the output interface is set to the opposite value.
If the -f is missing, both snmp interface indexes will be set to 0. [Default]

To learn more about Netflow version 5 datagram formats, see the [Cisco Netflow documentation](http://www.cisco.com/c/en/us/td/docs/net_mgmt/netflow_collection_engine/3-6/user/guide/format.html)
