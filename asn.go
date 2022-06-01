package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strconv"
	"strings"
)

// IP holds the BGP origin information about a given IP address.
type IP struct {
	ASNum     uint32 `json:"as_num"`
	IP        string `json:"ip"`
	BGPPrefix string `json:"bgp_prefix"`
	Country   string `json:"country"`
	Registry  string `json:"registry"`
	Allocated string `json:"allocated"`
	ASName    string `json:"as_name"`
}

// ASN holds the description of a BGP ASN.
type ASN struct {
	ASNum     uint32 `json:"as_num"`
	Country   string `json:"country"`
	Registry  string `json:"registry"`
	Allocated string `json:"allocated"`
	ASName    string `json:"as_name"`
}

const hexDigit = "0123456789abcdef"

func reverseaddr(addr string) (string, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", fmt.Errorf("unrecognized address: %s", addr)
	}
	if v4 := ip.To4(); v4 != nil {
		buf := make([]byte, 0, net.IPv4len*4)
		// Add it, in reverse, to the buffer
		for i := len(v4) - 1; i >= 0; i-- {
			buf = strconv.AppendInt(buf, int64(v4[i]), 10)
			// Only append a trailing "." if this isn't the final octet
			if i > 0 {
				buf = append(buf, '.')
			}
		}
		return string(buf), nil
	}

	buf := make([]byte, 0, net.IPv6len*4)
	for i := len(ip) - 1; i >= 0; i-- {
		v := ip[i]
		buf = append(buf, hexDigit[v&0xF])
		buf = append(buf, '.')
		buf = append(buf, hexDigit[v>>4])
		if i > 0 {
			buf = append(buf, '.')
		}
	}
	return string(buf), nil
}

// Parse the text output from the IP to ASN service and return an IP.
func parseOrigin(txt string) (IP, error) {
	fields := strings.Split(txt, "|")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	asn, err := strconv.ParseUint(fields[0], 10, 32)
	if err != nil && fields[0] != "NA" {
		return IP{}, errors.Wrap(err, "AS parsing failed")
	}

	return IP{
		ASNum:     uint32(asn),
		BGPPrefix: fields[1],
		Country:   fields[2],
		Registry:  fields[3],
		Allocated: fields[4],
	}, nil
}

var (
	originV4 = "origin.asn.cymru.com"
	originV6 = "origin6.asn.cymru.com"
)

// LookupIP queries Team Cymru's IP to ASN mapping service and returns BGP
// origin information about the IP.
func LookupIP(ip string) (IP, error) {
	rev, err := reverseaddr(ip)
	if err != nil {
		return IP{}, errors.Wrap(err, "reversing IP failed")
	}

	var zone string
	parsedIP := net.ParseIP(ip)
	if v4 := parsedIP.To4(); v4 != nil {
		zone = originV4
	} else {
		zone = originV6
	}

	q := fmt.Sprintf("%s.%s.", rev, zone)
	recs, err := net.LookupTXT(q)
	if err != nil {
		return IP{}, errors.Wrap(err, "DNS lookup failed")
	}

	origin, err := parseOrigin(recs[0])
	if err != nil {
		return IP{}, errors.Wrap(err, "parse failed")
	}
	origin.IP = ip
	if asn, err := LookupASN(fmt.Sprintf("AS%d", origin.ASNum)); err == nil {
		origin.ASName = asn.ASName
	}

	return origin, nil
}

// Parse the text output from the IP to ASN service and return an ASN.
func parseASN(txt string) (ASN, error) {
	fields := strings.Split(txt, "|")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}

	asn, err := strconv.ParseUint(fields[0], 10, 32)
	if err != nil && fields[0] != "NA" {
		return ASN{}, errors.Wrap(err, "AS parsing failed")
	}

	return ASN{
		ASNum:     uint32(asn),
		Country:   fields[1],
		Registry:  fields[2],
		Allocated: fields[3],
		ASName:    fields[4],
	}, nil
}

// LookupASN queries the IP to ASN service to fetch an AS description.
func LookupASN(asn string) (ASN, error) {
	if strings.ToLower(asn[0:2]) != "as" {
		asn = "AS" + asn
	}
	q := fmt.Sprintf("%s.asn.cymru.com.", asn)
	res, err := net.LookupTXT(q)
	if err != nil {
		return ASN{}, errors.Wrap(err, "DNS lookup failed")
	}
	as, err := parseASN(res[0])
	if err != nil {
		return ASN{}, errors.Wrap(err, "parse failed")
	}
	return as, nil
}
