## Usage

### Docker Image Run (Easiest)

Simply run in a container and pass any arguments at runtime. Below is an example passing the `--help` flag:

```
docker run -it --rm networkstatic/nflow_generator --help
```

To generate mock flow data simply add the target IP and port:

```
docker run -it --rm networkstatic/nflow_generator -t <ip> -p <port>
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

