// using for testing flow metrics for stats analysis
// also handy for integration testing
// Run using:
// go run nflow-generator.go nflow_logging.go nflow_payload.go  -t 172.16.86.138 -p 9995
// Or:
// go build
// ./nflow-generator -t <ip> -p <port>
package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	var opts struct {
		CollectorIP   string `short:"t" long:"target" description:"target ip address of the netflow collector"`
		CollectorPort string `short:"p" long:"port" description:"port number of the target netflow collector"`
		Help          bool   `short:"h" long:"help" description:"show nflow-generator help"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		showUsage()
		os.Exit(1)
	}
	if opts.Help == true {
		showUsage()
		os.Exit(1)
	}
	if opts.CollectorIP == "" || opts.CollectorPort == "" {
		showUsage()
		os.Exit(1)
	}
	collector := opts.CollectorIP + ":" + opts.CollectorPort
	udpAddr, err := net.ResolveUDPAddr("udp", collector)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Error connectiong to the target collector: ", err)
	}
	log.Infof("sending netflow data to a collector ip: %s and port: %s \n"+
		"use ctrl^c to terminate the app.", opts.CollectorIP, opts.CollectorPort)

	for {
		rand.Seed(time.Now().Unix())
		n := randomNum(50, 1000)
		data := GenerateNetflow(15)
		buffer := BuildNFlowPayload(data)
		_, err := conn.Write(buffer.Bytes())
		if err != nil {
			log.Fatal("Error connectiong to the target collector: ", err)
		}
		sleepInt := time.Duration(n)
		time.Sleep(sleepInt * time.Millisecond)
	}
}

func randomNum(min, max int) int {
	return rand.Intn(max-min) + min
}

func showUsage() {
	var usage string
	usage = `
Usage:
  main [OPTIONS] [collector IP address] [collector port number]

Application Options:
  -t, --target= target ip address of the netflow collector
  -p, --port=   port number of the target netflow collector

Example:
    ./nflow-generator.go -t 172.16.86.138 -p 9995

Help Options:
  -h, --help    Show this help message
  `
	fmt.Print(usage)
}
