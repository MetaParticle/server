package dns

const DNS_HEADER_LENGTH = 32

type DnsPacket struct {
	Header DnsHeader
}

type DnsHeader struct {
	Id, Meta, QDCount, ANCount, NSCount, ARCount int16
}

type DnsQuestion struct {
	QName string
	QType string
	QClass string
}