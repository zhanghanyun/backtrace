package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"net"
	"time"
)

type Result struct {
	i int
	s string
}

var rIp = []string{"219.141.136.12", "202.106.50.1", "221.179.155.161", "202.96.209.133", "210.22.97.1", "211.136.112.200", "58.60.188.222", "210.21.196.6", "120.196.165.24"}
var rName = []string{"北京电信", "北京联通", "北京移动", "上海电信", "上海联通", "上海移动", "广州电信", "广州联通", "广州移动"}
var ca = []color.Attribute{color.FgHiBlue, color.FgHiMagenta, color.FgHiYellow, color.FgHiGreen, color.FgHiCyan, color.FgHiRed, color.FgHiMagenta, color.FgHiYellow, color.FgHiBlue}
var m = map[string]string{"AS4134": "电信163 [普通线路]", "AS4809": "电信CN2 [优质线路]", "AS4837": "联通4837[普通线路]", "AS9929": "联通9929[优质线路]", "AS9808": "移动CMI [普通线路]", "AS58453": "移动CMI [普通线路]"}

func trace(ch chan Result, i int) {
	hops, err := Trace(net.ParseIP(rIp[i]))
	if err != nil {
		return
	}
	for _, h := range hops {
		for _, n := range h.Nodes {
			ip, err := lookupIpInfo(n.IP.String())
			if err != nil {
				log.Println(err)
				continue
			}
			if as, ok := m[ip.Org]; ok {
				c := color.New(ca[i]).Add(color.Bold).SprintFunc()
				s := fmt.Sprintf("%v %-15s %-23s", rName[i], rIp[i], c(as))
				ch <- Result{i, s}
				return
			}
		}
	}
}

func main() {

	var s = [9]string{}
	var c = make(chan Result)
	log.Println("正在测试三网回程路由...")
	for i := range rIp {
		go trace(c, i)
	}
loop:
	for i := 0; i < 9; i++ {
		select {
		case o := <-c:
			s[o.i] = o.s
		case <-time.After(time.Second * 10):
			break loop
		}
	}
	for _, s2 := range s {
		log.Println(s2)
	}
}
