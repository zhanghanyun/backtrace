package main

import (
	"log"
	"time"
)

func main() {

	var (
		s [9]string
		c = make(chan Result)
		t = time.After(time.Second * 10)
	)

	log.Println("正在测试三网回程路由...")
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
