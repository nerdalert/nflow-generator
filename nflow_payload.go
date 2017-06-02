package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

// Start time for this instance, used to compute sysUptime
var StartTime = time.Now().UnixNano()

// current sysUptime in msec - recalculated in CreateNFlowHeader()
var sysUptime uint32 = 0

// Counter of flow packets that have been sent
var flowSequence uint32 = 0

const (
	FTP_PORT        = 21
	SSH_PORT        = 22
	DNS_PORT        = 53
	HTTP_PORT       = 80
	HTTPS_PORT      = 443
	NTP_PORT        = 123
	SNMP_PORT       = 161
	IMAPS_PORT      = 993
	MYSQL_PORT      = 3306
	HTTPS_ALT_PORT  = 8080
	P2P_PORT        = 6681
	BITTORRENT_PORT = 6682
	UINT16_MAX      = 65535
	PAYLOAD_AVG_MD  = 1024
	PAYLOAD_AVG_SM  = 256
)

// struct data from fach
type NetflowHeader struct {
	Version        uint16
	FlowCount      uint16
	SysUptime      uint32
	UnixSec        uint32
	UnixMsec       uint32
	FlowSequence   uint32
	EngineType     uint8
	EngineId       uint8
	SampleInterval uint16
}

type NetflowPayload struct {
	SrcIP          uint32
	DstIP          uint32
	NextHopIP      uint32
	SnmpInIndex    uint16
	SnmpOutIndex   uint16
	NumPackets     uint32
	NumOctets      uint32
	SysUptimeStart uint32
	SysUptimeEnd   uint32
	SrcPort        uint16
	DstPort        uint16
	Padding1       uint8
	TcpFlags       uint8
	IpProtocol     uint8
	IpTos          uint8
	SrcAsNumber    uint16
	DstAsNumber    uint16
	SrcPrefixMask  uint8
	DstPrefixMask  uint8
	Padding2       uint16
}

//Complete netflow records
type Netflow struct {
	Header  NetflowHeader
	Records []NetflowPayload
}

//Marshall NetflowData into a buffer
func BuildNFlowPayload(data Netflow) bytes.Buffer {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, &data.Header)
	if err != nil {
		log.Println("Writing netflow header failed:", err)
	}
	for _, record := range data.Records {
		err := binary.Write(buffer, binary.BigEndian, &record)
		if err != nil {
			log.Println("Writing netflow record failed:", err)
		}
	}
	return *buffer
}

//Generate a netflow packet w/ user-defined record count
func GenerateNetflow(recordCount int) Netflow {
	data := new(Netflow)
	header := CreateNFlowHeader(recordCount)
	records := []NetflowPayload{}
	if recordCount == 8 {
		// overwrite payload to add some variations for traffic spikes.
		records = CreateVariablePayload(recordCount)
	} else {
		records = CreateNFlowPayload(recordCount)
	}

	data.Header = header
	data.Records = records
	return *data
}

//Generate and initialize netflow header
func CreateNFlowHeader(recordCount int) NetflowHeader {

	t := time.Now().UnixNano()
	sec := t / int64(time.Second)
	nsec := t - sec*int64(time.Second)
	sysUptime = uint32((t-StartTime) / int64(time.Millisecond))+1000
	flowSequence++

	// log.Infof("Time: %d; Seconds: %d; Nanoseconds: %d\n", t, sec, nsec)
	// log.Infof("StartTime: %d; sysUptime: %d", StartTime, sysUptime)
	// log.Infof("FlowSequence %d", flowSequence)

	h := new(NetflowHeader)
	h.Version = 5
	h.FlowCount = uint16(recordCount)
	h.SysUptime = sysUptime
	h.UnixSec = uint32(sec)
	h.UnixMsec = uint32(nsec)
	h.FlowSequence = flowSequence
	h.EngineType = 1
	h.EngineId = 0
	h.SampleInterval = 0
	return *h
}

func CreateVariablePayload(recordCount int) []NetflowPayload {
	payload := make([]NetflowPayload, recordCount)

	for i := 0; i < recordCount; i++ {
		payload[0] = CreateHttpFlow()
		payload[1] = CreateHttpsFlow()
		payload[2] = CreateHttpAltFlow()
		payload[3] = CreateDnsFlow()
		payload[5] = CreateNtpFlow()
		payload[6] = CreateImapsFlow()
		payload[7] = CreateMySqlFlow()
	}

	return payload
}

func CreateNFlowPayload(recordCount int) []NetflowPayload {
	payload := make([]NetflowPayload, recordCount)
	for i := 0; i < recordCount; i++ {
		payload[0] = CreateHttpFlow()
		payload[1] = CreateHttpsFlow()
		payload[2] = CreateHttpAltFlow()
		payload[3] = CreateDnsFlow()
		payload[4] = CreateIcmpFlow()
		payload[5] = CreateNtpFlow()
		payload[6] = CreateImapsFlow()
		payload[7] = CreateMySqlFlow()
		payload[8] = CreateRandomFlow()
		payload[9] = CreateSshFlow()
		payload[10] = CreateP2pFlow()
		payload[11] = CreateBitorrentFlow()
		payload[12] = CreateFTPFlow()
		payload[13] = CreateSnmpFlow()
		payload[14] = CreateIcmpFlow()
		payload[15] = CreateRandomFlow()
	}
	return payload
}

//Initialize netflow record with random data
func CreateIcmpFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("172.16.50.10")
	payload.DstIP = IPtoUint32("132.12.130.10")
	payload.NextHopIP = IPtoUint32("132.12.130.1")
	payload.SrcPort = 0
	payload.DstPort = 0
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_SM)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_SM)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 1
	// payload.IpTos = 0
	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_SM, 1, rand.Intn(32))
	return *payload
}

//Initialize netflow record with random data
func CreateHttpFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("112.10.20.10")
	payload.DstIP = IPtoUint32("172.30.190.10")
	payload.NextHopIP = IPtoUint32("172.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(HTTP_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

//Initialize netflow record with random data
func CreateSnmpFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("112.10.20.10")
	payload.DstIP = IPtoUint32("172.30.190.10")
	payload.NextHopIP = IPtoUint32("172.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(SNMP_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 17
	// payload.IpTos = 0
	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 17, rand.Intn(32))
	return *payload
}

func CreateFTPFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("112.10.100.10")
	payload.DstIP = IPtoUint32("192.168.120.10")
	payload.NextHopIP = IPtoUint32("172.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(FTP_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateNtpFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("247.104.20.202")
	payload.DstIP = IPtoUint32("10.12.190.10")
	payload.NextHopIP = IPtoUint32("192.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(NTP_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 17
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)
	// payload.SrcPrefixMask = uint8(32)
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 17, 32)
	return *payload
}

func CreateP2pFlow() NetflowPayload {
	payload := new(NetflowPayload)
	payload.SrcIP = IPtoUint32("247.104.20.202")
	payload.DstIP = IPtoUint32("10.12.190.10")
	payload.NextHopIP = IPtoUint32("192.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(P2P_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 17
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(32)
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 17, 32)
	return *payload
}

func CreateBitorrentFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("192.168.20.202")
	payload.DstIP = IPtoUint32("42.12.190.10")
	payload.NextHopIP = IPtoUint32("192.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(BITTORRENT_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 17
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(32)
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 17, 32)
	return *payload
}

func CreateSshFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("172.30.20.102")
	payload.DstIP = IPtoUint32("222.12.190.10")
	payload.NextHopIP = IPtoUint32("192.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(SSH_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateHttpsFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("192.168.20.10")
	payload.DstIP = IPtoUint32("202.12.190.10")
	payload.NextHopIP = IPtoUint32("172.199.15.1")
	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(HTTPS_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateHttpAltFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("10.10.20.122")
	payload.DstIP = IPtoUint32("84.12.190.210")
	payload.NextHopIP = IPtoUint32("192.199.15.1")
	payload.SrcPort = uint16(12001)
	payload.DstPort = uint16(HTTPS_ALT_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateDnsFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("59.220.158.122")
	payload.DstIP = IPtoUint32("10.12.233.210")
	payload.NextHopIP = IPtoUint32("39.199.15.1")
	payload.SrcPort = uint16(9221)
	payload.DstPort = uint16(DNS_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 17
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 17, rand.Intn(32))
	return *payload
}

func CreateImapsFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("172.30.20.102")
	payload.DstIP = IPtoUint32("62.12.190.10")
	payload.NextHopIP = IPtoUint32("131.199.15.1")
	payload.SrcPort = uint16(9010)
	payload.DstPort = uint16(IMAPS_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateMySqlFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = IPtoUint32("10.154.20.12")
	payload.DstIP = IPtoUint32("77.12.190.94")
	payload.NextHopIP = IPtoUint32("150.20.145.1")
	payload.SrcPort = uint16(9010)
	payload.DstPort = uint16(MYSQL_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

func CreateRandomFlow() NetflowPayload {
	payload := new(NetflowPayload)

	payload.SrcIP = rand.Uint32()
	payload.DstIP = rand.Uint32()
	payload.NextHopIP = rand.Uint32()
	payload.SrcPort = genRandUint16(UINT16_MAX)
	payload.DstPort = genRandUint16(UINT16_MAX)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	// payload.NumPackets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.NumOctets = genRandUint32(PAYLOAD_AVG_MD)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	// payload.Padding1 = 0
	// payload.IpProtocol = 6
	// payload.IpTos = 0
	// payload.SrcAsNumber = genRandUint16(UINT16_MAX)

	// payload.SrcPrefixMask = uint8(rand.Intn(32))
	// payload.DstPrefixMask = uint8(rand.Intn(32))
	// payload.Padding2 = 0
	FillCommonFields(payload, PAYLOAD_AVG_MD, 6, rand.Intn(32))
	return *payload
}

// patch up the common fields of the packets
func FillCommonFields (
		payload *NetflowPayload, 
		numPktOct int, 
		ipProtocol int, 
		srcPrefixMask int) NetflowPayload {

// Fill template with values not filled by caller
	// payload.SrcIP = IPtoUint32("10.154.20.12")
	// payload.DstIP = IPtoUint32("77.12.190.94")
	// payload.NextHopIP = IPtoUint32("150.20.145.1")
	// payload.SrcPort = uint16(9010)
	// payload.DstPort = uint16(MYSQL_PORT)
	// payload.SnmpInIndex = genRandUint16(UINT16_MAX)
	// payload.SnmpOutIndex = genRandUint16(UINT16_MAX)
	payload.NumPackets = genRandUint32(numPktOct)
	payload.NumOctets = genRandUint32(numPktOct)
	// payload.SysUptimeStart = rand.Uint32()
	// payload.SysUptimeEnd = rand.Uint32()
	payload.Padding1 = 0
	payload.IpProtocol = uint8(ipProtocol)
	payload.IpTos = 0
	payload.SrcAsNumber = genRandUint16(UINT16_MAX)
	payload.DstAsNumber = genRandUint16(UINT16_MAX)

	payload.SrcPrefixMask = uint8(srcPrefixMask)
	payload.DstPrefixMask = uint8(rand.Intn(32))
	payload.Padding2 = 0

	// now handle computed values
	if !opts.FalseIndex {                       // default interfaces are zero
		payload.SnmpInIndex = 0
		payload.SnmpOutIndex = 0
	} else if payload.SrcIP > payload.DstIP {   // false-index
		payload.SnmpInIndex = 1
		payload.SnmpOutIndex = 2
	} else {
		payload.SnmpInIndex = 2
		payload.SnmpOutIndex = 1
	}

	uptime := int(sysUptime)
	payload.SysUptimeEnd = uint32(uptime - randomNum(10,500))
	payload.SysUptimeStart = payload.SysUptimeEnd - uint32(randomNum(10,500))

	// log.Infof("S&D : %x %x %d, %d", payload.SrcIP, payload.DstIP, payload.DstPort, payload.SnmpInIndex)
	// log.Infof("Time: %d %d %d", sysUptime, payload.SysUptimeStart, payload.SysUptimeEnd)

	return *payload
}

func genRandUint16(max int) uint16 {
	return uint16(rand.Intn(max))
}

func IPtoUint32(s string) uint32 {
	ip := net.ParseIP(s)
	return binary.BigEndian.Uint32(ip.To4())
}

func genRandUint32(max int) uint32 {
	return uint32(rand.Intn(max))
}
