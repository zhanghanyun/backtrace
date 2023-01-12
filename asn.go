package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"strings"
)

type Result struct {
	i int
	s string
}

var (
	rIp   = []string{"219.141.136.12", "202.106.50.1", "221.179.155.161", "202.96.209.133", "210.22.97.1", "211.136.112.200", "58.60.188.222", "210.21.196.6", "120.196.165.24", "61.139.2.69", "119.6.6.6", "211.137.96.205"}
	rName = []string{"北京电信", "北京联通", "北京移动", "上海电信", "上海联通", "上海移动", "广州电信", "广州联通", "广州移动", "成都电信", "成都联通", "成都移动"}
	ca    = []color.Attribute{color.FgHiYellow, color.FgHiMagenta, color.FgHiBlue, color.FgHiGreen, color.FgHiCyan, color.FgHiRed}
	m     = map[string]string{"AS4134": "电信163 [普通线路]", "AS4809": "电信CN2 [优质线路]", "AS4837": "联通4837[普通线路]", "AS9929": "联通9929[优质线路]", "AS9808": "移动CMI [普通线路]", "AS58453": "移动CMI [普通线路]"}
)

func trace(ch chan Result, i int) {

	hops, err := Trace(net.ParseIP(rIp[i]))
	if err != nil {
		s := fmt.Sprintf("%v %-15s %v", rName[i], rIp[i], err)
		ch <- Result{i, s}
		return
	}

	for _, h := range hops {
		for _, n := range h.Nodes {
			asn := ipAsn(n.IP.String())
			if asn == "" {
				continue
			} else {
				as := m[asn]
				cai := i % 6
				c := color.New(ca[cai]).Add(color.Bold).SprintFunc()
				s := fmt.Sprintf("%v %-15s %-23s", rName[i], rIp[i], c(as))
				ch <- Result{i, s}
				return
			}
		}
	}

	s := fmt.Sprintf("%v %-15s %v", rName[i], rIp[i], "测试超时")
	ch <- Result{i, s}
}

func ipAsn(ip string) string {

	switch {
	case strings.HasPrefix(ip, "59.43"):
		return "AS4809"
	case strings.HasPrefix(ip, "202.97"):
		return "AS4134"
	case strings.HasPrefix(ip, "218.105") || strings.HasPrefix(ip, "210.51"):
		return "AS9929"
	case strings.HasPrefix(ip, "219.158"):
		return "AS4837"
	case strings.HasPrefix(ip, "223.118") || strings.HasPrefix(ip, "223.119") || strings.HasPrefix(ip, "223.120") || strings.HasPrefix(ip, "223.121"):
		return "AS58453"
	default:
		return ""
	}
}
