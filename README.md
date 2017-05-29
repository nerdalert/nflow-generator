## Usage

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

Install [Go](http://golang.org/doc/install)

	git clone https://github.com/nerdalert/nflow-generator.git
	cd <dir>
	go build

Go build will leave a binary in the root directory that can be run.
	
### RUN

Feed it the target collector and port:

	./nflow-generator -t <ip> -p <port>

Or:

	go run nflow-generator.go nflow_logger.go nflow_data.go  -t 172.16.86.138 -p 9995

### Update - May 2017

The original mock netflow data set random values in several fields which confused certain netflow collectors.
Those collectors complained about inaccurate time stamps, 
and were confused by the random values sent in the input and output interface fields. This update:

* Sets the SysUptime, unix_secs, and unix_nsecs fields to sensible (UTC) values
* Generates a unique flow_sequence value for each netflow datagram
* Creates reasonable start/stop times for flows, so the First is (now-X) and Last is (now-Y), where X & Y are random times, and X > Y.
* Sets the interface indexes to 1 or 2 - based on this algorithm. 
If the source address > dest address, input interface is set to 1, and set to 2 otherwise,
and the output interface is set to the opposite value.

These are based on the [Cisco Netflow documentation](http://www.cisco.com/c/en/us/td/docs/net_mgmt/netflow_collection_engine/3-6/user/guide/format.html)