## Usage

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

