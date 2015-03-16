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

type Proto int

const (
	FTP Proto = iota + 1
	SSH
	DNS
	HTTP
	HTTPS
	NTP
	SNMP
	IMAPS
	MYSQL
	HTTPS_ALT
	P2P
	BITTORRENT
)

var opts struct {
	CollectorIP   string `short:"t" long:"target" description:"target ip address of the netflow collector"`
	CollectorPort string `short:"p" long:"port" description:"port number of the target netflow collector"`
	SpikeProto    string `short:"s" long:"spike" description:"run a second thread generating a spike for the specified protocol"`
	Help          bool   `short:"h" long:"help" description:"show nflow-generator help"`
}

func main() {

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
	log.Infof("sending netflow data to a collector ip: %s and port: %s. \n"+
		"Use ctrl^c to terminate the app.", opts.CollectorIP, opts.CollectorPort)

	for {
		rand.Seed(time.Now().Unix())
		n := randomNum(50, 1000)
		// add spike data
		if opts.SpikeProto != "" {
			GenerateSpike()
		}
		if n > 900 {
			data := GenerateNetflow(8)
			buffer := BuildNFlowPayload(data)
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				log.Fatal("Error connectiong to the target collector: ", err)
			}
		} else {
			data := GenerateNetflow(16)
			buffer := BuildNFlowPayload(data)
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				log.Fatal("Error connectiong to the target collector: ", err)
			}
		}
		// add some periodic spike data
		if n < 150 {
			sleepInt := time.Duration(3000)
			time.Sleep(sleepInt * time.Millisecond)
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
  -s, --spike run a second thread generating a spike for the specified protocol
    protocol options are as follows:
        ftp - generates tcp/21
        ssh  - generates tcp/22
        dns - generates udp/54
        http - generates tcp/80
        https - generates tcp/443
        ntp - generates udp/123
        snmp - generates ufp/161
        imaps - generates tcp/993
        mysql - generates tcp/3306
        https_alt - generates tcp/8080
        p2p - generates udp/6681
        bittorrent - generates udp/6682

Example:
    -generate default flows:
    ./nflow-generator.go -t 172.16.86.138 -p 9995

    -generate default flows along with a spike in the specified protocol:
    ./nflow-generator -t 172.16.86.138 -p 9995 -s ssh

Help Options:
  -h, --help    Show this help message
  `
	fmt.Print(usage)
}
