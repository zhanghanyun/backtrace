package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"time"
)

type IpInfo struct {
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Org     string `json:"org"`
}

func main() {

	var (
		s [12]string
		c = make(chan Result)
		t = time.After(time.Second * 10)
	)

	go func() {
		http.Get("https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fzhanghanyun%2Fbacktrace&count_bg=%2379C83D&title_bg=%23555555&icon=&icon_color=%23E7E7E7&title=hits&edge_flat=false")
	}()

	yellow := color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	log.Println("正在测试三网回程路由...")

	rsp, _ := http.Get("http://ipinfo.io")
	info := IpInfo{}
	json.NewDecoder(rsp.Body).Decode(&info)

	fmt.Println(green("国家: ") + cyan(info.Country) + green(" 城市: ") + cyan(info.City) + green(" 服务商: ") + cyan(info.Org))
	fmt.Println(green("项目地址:"), yellow("https://github.com/zhanghanyun/backtrace"))

	for i := range ips {
		go trace(c, i)
	}

loop:
	for range s {
		select {
		case o := <-c:
			s[o.i] = o.s
		case <-t:
			break loop
		}
	}

	for _, r := range s {
		fmt.Println(r)
	}
	log.Println(green("测试完成!"))
}
