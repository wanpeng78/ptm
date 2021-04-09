package filter

import (
	"log"
	"sync"
	"testing"
)

func TestSpeedTest(t *testing.T) {
	var wg sync.WaitGroup
	urls := []string{
		"mirrors.tuna.tsinghua.edu.cn",
		"https://www.baidu.com",
		"www.google.com",
		"www.da.cc",
		"123.22.33.2",
		"https://www.baidu.com/home",
	}
	c := make(chan int, 1)
	go func() {
		for {
			select {
			case idx := <-c:
				log.Println(urls[idx])
			}
		}
	}()
	SpeedTest(urls, c)
	wg.Wait()
}
