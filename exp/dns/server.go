package dns

/*
-- INCOMPLETE--
This is an implementation of a DNS name server written in Go.
It uses TCP instead of UDP, since it is a part of the specification and I do not want to deal with a situation where I have to switch to TCP just becuase a package is to big.
Might switch back if I see that packages don't grow to big.
It is written after the RFC 1035 specification that can be found at http://www.freesoft.org/CIE/RFC/1035/.
It handles A, NS and TXT types and the IN class.
No inverse querys at the moment.
No recursive querys at the moment.
If anyone wishes to implement MX and other types, go ahead (no pun intended)!
*/

import (
	//"github.com/MetaParticle/metaparticle/logger"
	
	"net"
	"bytes"
)

func ListenAndServe(addr string) {
	listener, err := net.ListenTCP("tcp", addr)
	for err == nil {
		conn, ConnErr := listener.Accept()
		if ConnErr != nil {
			go handleConn(conn)
		}
	}
}

var bb = make(chan []byte, 100)

func readPacket(conn *net.Conn) ([]byte, error) {
	var size int
	for size, err := readSize(conn); err != nil {
		size, err = readSize(conn)
	}
	b := make([]byte, size, 1024)
	n, err := conn.Read(b)
	return b[:n], err
}

func readSize(conn *net.Conn) (int, error) {
	sizeSlice := make([]byte, 2, 2)
	n, err := conn.Read(sizeSlice)
	size := int16(sizeSlice[0])<<8 | int16(sizeSlice[1])
	return size, err
}

func handleConn(conn *net.Conn) {
	var b []byte
	for {
		for b, err := readPacket(conn); err != nil { //read until packet is found.
			b, err = readPacket(conn)
		}
		header := readHeader(b) //parse header for future use
		qs := readQuestions(header, b) //if error here, then "RCODE" = "Format Error"
		
		
		
		answers := 0
		nameservers := 0
		rcode := 5 //"Refused", becuase that is cooler than "Not Implemented".
		for _, q := range qs {
			//lookup answer to question
			//make answer section if have answer
			//else if have nameserver, make such.
			//else ignore and remember to set RCODE = "Name Error"
			
		}
		fixHeader(true, rcode, answers, nameservers, b) // prepare header for response
		ans := append(b, ansSlice)
		conn.Write(ans)
	}
}

func readHeader(b []byte) DnsHeader {
	buf := bytes.NewBuffer(b)
	h := new(DnsHeader)
	h.Id = buf.Read()<<8 | buf.Read()
	h.Meta = buf.Read()<<8 | buf.Read()
	h.QDCount = buf.Read()<<8 | buf.Read()
	h.ANCount = buf.Read()<<8 | buf.Read()
	h.NSCount = buf.Read()<<8 | buf.Read()
	h.ARCount = buf.Read()<<8 | buf.Read()
	
	return h
}

/*
This appears to be the hard part, as I have no clue of how effectively 
parse to labels in the "QNAME" field.
*/
func readQuestions(header DnsHeader, b []byte) []*DnsQuestion {
	qs := make([]DnsQuestion, 0, header.QDCount)
	off := DNS_HEADER_LENGTH
	for i := 0; i < header.QDCount; i++ {
		ls, off := parseLabels(b, off)
		t := string(b[off]) + string(b[off+1])
		c := string(b[off+2]) + string(b[off+3])
		off += 4
		for name := range ls {
			qs = append(qs, &DnsQuestion{name, t, c})
		}
	}
}

func parseLabels(msg []byte, start int) []string, int {
	offset := start
	ok := true
	res := make([]string, 0, 10)
	for ok {
		s, offset, ok := unpackDomainName(msg, offset)
		if ok {
			res := append(res, s)
		}
	}
	return res, offset
}

/*
Used to set the header bits in the response.
*/
func fixHeader(aa bool, rcode byte, an int16, ns int16, b []byte) {
	b[2] = b[2] | 0b10000000
	if aa {
		b[2] = b[2] | 0b00000100//if we are authoritative, which we should be
	}
	b[3] = b[3] & 0 //no recursion, Z must be 0 and RCODE is below.
	b[3] = b[3] | (rcode & 7) //set responsecode
	
	b[6], b[7] = byte(an >> 8), byte(an) //number of answers
	b[8], b[9] = byte(ns >> 8), byte(ns) //number if nameservers
	b[10], b[11] = 0, 0 //just to make sure.
}

func


/*
 This part is "borrowed" from "dnsmsg.go" in the "net" package of standard Go.
 Since that code is supposed to be BSD-like licensed, I can't see any licensing problems.
 It is already in the code-base, only not public.
 All I want is accsses to it
 
 After seeing it, I might as well use it since I would probably write something very like it on my own after that.
 */
func unpackDomainName(msg []byte, off int) (s string, off1 int, ok bool) {
	s = ""
	ptr := 0 // number of pointers followed
	Loop:
		for {
			if off >= len(msg) {
					return "", len(msg), false
			}
			c := int(msg[off])
			off++
			switch c & 0xC0 {
			case 0x00:
				if c == 0x00 {
					// end of name
					break Loop
				}
				// literal string
				if off+c > len(msg) {
						return "", len(msg), false
				}
				s += string(msg[off:off+c]) + "."
				off += c
			case 0xC0:
				// pointer to somewhere else in msg.
				// remember location after first ptr,
				// since that's how many bytes we consumed.
				// also, don't follow too many pointers --
				// maybe there's a loop.
				if off >= len(msg) {
					return "", len(msg), false
				}
				c1 := msg[off]
				off++
				if ptr == 0 {
					off1 = off
				}
				if ptr++; ptr > 10 {
					return "", len(msg), false
				}
				off = (c^0xC0)<<8 | int(c1)
			default:
				// 0x80 and 0x40 are reserved
				return "", len(msg), false
				}
		}
		if ptr == 0 {
			off1 = off
		}
		return s, off1, true
}

