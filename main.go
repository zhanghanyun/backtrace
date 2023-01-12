package main

import (
	"github.com/fatih/color"
	"log"
	"time"
)

func main() {

	var (
		s [12]string
		c = make(chan Result)
		t = time.After(time.Second * 10)
	)

	head := color.New(color.FgHiBlue).Add(color.Bold).SprintFunc()
	note := color.New(color.FgGreen).SprintFunc()
	log.Println(head("项目地址：github.com/zhanghanyun/backtrace"))
	log.Println(note("正在测试三网回程路由..."))

	for i := range rIp {
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
		log.Println(r)
	}
}
