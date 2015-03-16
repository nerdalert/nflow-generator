package main

//import "fmt"

//Generate a netflow packet w/ user-defined record count
func GenerateSpike() Netflow {
	data := new(Netflow)
	header := CreateNFlowHeader(1)
	records := []NetflowPayload{}
	records = spikeFlowPayload()
	data.Header = header
	data.Records = records
	return *data
}

func spikeFlowPayload() []NetflowPayload {
	payload := make([]NetflowPayload, 1)
	switch opts.SpikeProto {
	case "ssh":
		payload[0] = CreateSshFlow()
	case "ftp":
		payload[0] = CreateFTPFlow()
	case "http":
		payload[0] = CreateHttpFlow()
	case "https":
		payload[0] = CreateHttpsFlow()
	case "ntp":
		payload[0] = CreateNtpFlow()
	case "snmp":
		payload[0] = CreateSnmpFlow()
	case "imaps":
		payload[0] = CreateImapsFlow()
	case "mysql":
		payload[0] = CreateMySqlFlow()
	case "https_alt":
		payload[0] = CreateHttpAltFlow()
	case "p2p":
		payload[0] = CreateP2pFlow()
	case "bittorrent":
		payload[0] = CreateBitorrentFlow()
	default:
		log.Fatalf("protocol option %s is not valid, see --help for options", opts.SpikeProto)
	}
	return payload
}
