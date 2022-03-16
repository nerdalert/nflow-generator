# Usage - nflow-generator

[![nflow-generator image CI](https://github.com/nerdalert/nflow-generator/actions/workflows/image-build.yml/badge.svg)](https://github.com/nerdalert/nflow-generator/actions/workflows/image-build.yml)

This program generates mock netflow (v5) data that can be used to test netflow collector programs. 
The program simulates a router that is exporting flow records to the collector.
It is useful for determining whether the netflow collector is operating and/or receiving netflow datagrams.

nflow-generator generates several netflow datagrams per second, each with 8 or 16 records for varying kinds of traffic (HTTP, SSH, SNMP, DNS, MySQL, and many others.)

### Docker Image Run (Easiest)

Simply run in a container and pass any arguments at runtime. Below is an example passing the `--help` flag:

```
docker run -it --rm networkstatic/nflow-generator --help
# or podman/quay repos
podman run -it --rm /quay.io/networkstatic/nflow-generator --help
```

To generate mock flow data simply add the target IP and port:

```
docker run -it --rm networkstatic/nflow-generator -t <ip> -p <port>
# or podman/quay repos
podman run -it --rm /quay.io/networkstatic/nflow-generator -t <ip> -p <port>
```

### Download the binary

You can download the Linux binary here [nflow-generator-x86_64-linux](https://github.com/nerdalert/nflow-generator/blob/master/binaries/nflow-generator-x86_64-linux).
### Build

Install [Go](http://golang.org/doc/install), then:

	git clone https://github.com/nerdalert/nflow-generator.git 
	cd nflow-generator
	go build

Go build will leave a binary in the root directory that can be run.
	
### RUN

Feed it the target collector and port, and optional "false-index" flag:

	./nflow-generator -t <ip> -p <port> [ -f | --false-index ]

### Run a Test Collection

You can run a simple test collection using nfcapd from the nfdump package with the following.

- Start a netflow collector

```
sudo apt-get install nfdump
mkdir /tmp/nfcap-test
nfcapd -E  -p 9001 -l /tmp/nfcap-test
```

In a seperate console, run the netflow-generator pointing at an IP on the host the collector is running on (in this case the VM has an IP of 192.168.1.113).

```
sudo docker run -it --rm networkstatic/nflow-generator -t 192.168.1.113 -p 9001
```

- You should start seeing records displayed to the output of the screen running nfcapd like the following.

```
$> nfcapd -E  -p 9001 -l /tmp/nfcap-test
Add extension: 2 byte input/output interface index
Add extension: 4 byte input/output interface index
Add extension: 2 byte src/dst AS number
Add extension: 4 byte src/dst AS number
Bound to IPv4 host/IP: any, Port: 9001
Startup.
Init IPFIX: Max number of IPFIX tags: 62

Flow Record:
  Flags        =              0x00 FLOW, Unsampled
  export sysid =                 1
  size         =                56
  first        =        1552592037 [2019-03-14 15:33:57]
  last         =        1552592038 [2019-03-14 15:33:58]
  msec_first   =               973
  msec_last    =               414
  src addr     =      112.10.20.10
  dst addr     =     172.30.190.10
  src port     =                40
  dst port     =                80
  fwd status   =                 0
  tcp flags    =              0x00 ......
  proto        =                 6 TCP
  (src)tos     =                 0
  (in)packets  =               792
  (in)bytes    =                23
  input        =                 0
  output       =                 0
  src as       =             48730
  dst as       =             15401


Flow Record:
  Flags        =              0x00 FLOW, Unsampled
  export sysid =                 1
  size         =                56
  first        =        1552592038 [2019-03-14 15:33:58]
  last         =        1552592038 [2019-03-14 15:33:58]
  msec_first   =               229
  msec_last    =               379
  src addr     =     192.168.20.10
  dst addr     =     202.12.190.10
  src port     =                40
  dst port     =               443
  fwd status   =                 0
  tcp flags    =              0x00 ......
  proto        =                 6 TCP
  (src)tos     =                 0
  (in)packets  =               599
  (in)bytes    =               602
  input        =                 0
  output       =                 0
  src as       =              1115
  dst as       =             50617

```

### Notes

The original mock netflow generator placed random values in several fields which confused 
certain netflow collectors that complained about inaccurate time stamps, 
and were confused by the random values sent in the input and output interface fields. 

Changes:

* Sets the `SysUptime`, `unix_secs`, and `unix_nsecs` fields of the Netflow datagrams to sensible (UTC) values
* Generates a unique `flow_sequence` value for each netflow datagram
* Creates reasonable start/stop times for flows, so the First is set to (now-X) and Last to (now-Y), where X & Y are random times, and X > Y.
* If the --false-index (-f) flag is set on the command line, 
use this algorithm to set the interface indexes to 1 or 2:
If the source address > dest address, input interface is set to 1, and set to 2 otherwise,
and the output interface is set to the opposite value.
If the -f is missing, both snmp interface indexes will be set to 0. [Default]

To learn more about Netflow version 5 datagram formats, see the [Cisco Netflow documentation](http://www.cisco.com/c/en/us/td/docs/net_mgmt/netflow_collection_engine/3-6/user/guide/format.html)
