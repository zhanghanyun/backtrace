package main

import (
	"github.com/fatih/color"
	"log"
	"net"
	"sync"
)

var rIp = []string{"219.141.136.12", "202.106.50.1", "221.179.155.161", "202.96.209.133", "210.22.97.1", "211.136.112.200", "58.60.188.222", "210.21.196.6", "120.196.165.24"}
var rName = []string{"北京电信", "北京联通", "北京移动", "上海电信", "上海联通", "上海移动", "广州电信", "广州联通", "广州移动"}
var ca = []color.Attribute{color.FgHiBlue, color.FgHiMagenta, color.FgHiYellow, color.FgHiGreen, color.FgHiCyan, color.FgHiRed, color.FgHiMagenta, color.FgHiYellow, color.FgHiBlue}
var m = map[uint32]string{4134: "电信163 [普通线路]", 4809: "电信CN2 [优质线路]", 4837: "联通4837[普通线路]", 9929: "联通9929[优质线路]", 9808: "移动CMI [普通线路]", 58453: "移动CMI [普通线路]"}

func trace(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	hops, err := Trace(net.ParseIP(rIp[i]))
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range hops {
		for _, n := range h.Nodes {
			ip, err := LookupIP(n.IP.String())
			if err != nil {
				log.Fatal(err)
			}
			if ip.Country == "CN" {
				c := color.New(ca[i]).Add(color.Bold).SprintFunc()
				log.Printf("%v %-15s %-23s %dms\n", rName[i], rIp[i], c(m[ip.ASNum]), n.RTT[0].Milliseconds())
				return
			}
		}
	}
}

func main() {

	var wg = sync.WaitGroup{}
	log.Println("正在测试三网回程路由...")
	for i := range rIp {
		wg.Add(1)
		go trace(&wg, i)
	}
	wg.Wait()
}
