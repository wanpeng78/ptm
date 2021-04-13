package filter

import (
	"errors"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// Result is just result
type Result struct {
	//Address is IP or domain
	Address string
	//Latency in ms
	Latency time.Duration
	//Error for invalid remote host
	Error error
}

// Results for results
type Results []Result

// SortByDuration just sort follow duration
// Using insert sort
func (r *Results) SortByDuration() {
	re := insertSort(*r)
	r = &re
}

// LatencyTestWithChan test latency between remote and local,
// and store valid remote's index in c
// the faster the less for remote's index
//
// Args:
// addrs -> []string{"www.bing.com","192.168.12.1"}
// c 	 -> store result for goroutine
// loops -> ping times
// retry -> wheather try to use net.Dial to test latency after ping failed
// timeout -> timeout for ping and net.Dial
//
// this is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTestWithChan(addrs []string, c chan<- Result, loops int, retry bool, timeout time.Duration) error {
	if len(addrs) < 1 {
		return errors.New("remote host addr must more than one")
	}
	// Check privileged
	if os.Geteuid() != 0 {
		return errors.New("please run this program as root user")
	}
	fixedAddrs := urlHandler(addrs)
	var wg sync.WaitGroup

	for idx, addr := range fixedAddrs {
		wg.Add(1)
		go func(addr string, idx int) {
			defer wg.Done()
			var err error
			var latency time.Duration
			// By using icmp(ping) to test latency
			p := ping.New(addr)
			p.Timeout = timeout
			p.Count = loops
			err = p.Run()
			statis := p.Statistics()
			// handle remote server banned icmp(ping)
			// switch tp tcp:http
			if statis.PacketLoss == 100 && retry {
				last := time.Now()
				_, err = net.DialTimeout("tcp", addr+":http", timeout)
				latency = time.Now().Sub(last)
			} else {
				latency = statis.AvgRtt
			}
			c <- Result{Address: addr, Latency: latency, Error: err}
		}(addr, idx)
	}
	wg.Wait()
	return nil
}

// LatencyTest test latency between remote and local,
// and store valid remote's index in c
// tresult are stored in rets,the faster remote host the less index are
//
// Args:
// addrs -> []string{"www.bing.com","192.168.12.1"}
// loops -> ping times
// retry -> wheather try to use net.Dial to test latency after ping failed
// timeout -> timeout for ping and net.Dial
//
// for improving efficiency to use LatencyTestWithChan in goroutine instead of this
// but this provide a sort method(Results.SortByDuration)
// This is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTest(addrs []string, loops int, retry bool, timeout time.Duration) (rets Results, err error) {
	if len(addrs) < 1 {
		return nil, errors.New("remote host addr must more than one")
	}
	// Check privileged
	if os.Geteuid() != 0 {
		return nil, errors.New("please run this program as root user")
	}
	fixedAddrs := urlHandler(addrs)
	var wg sync.WaitGroup

	for idx, addr := range fixedAddrs {
		wg.Add(1)
		go func(addr string, idx int) {
			defer wg.Done()
			var err error
			var latency time.Duration
			p := ping.New(addr)
			p.Timeout = timeout
			p.Count = loops
			err = p.Run()
			statis := p.Statistics()
			// handle remote server banned icmp(ping)
			// switch tp tcp:http
			if statis.PacketLoss == 100 && retry {
				last := time.Now()
				_, err = net.DialTimeout("tcp", addr+":http", timeout)
				latency = time.Now().Sub(last)
			} else {
				latency = statis.AvgRtt
			}
			rets = append(rets, Result{Address: addr, Latency: latency, Error: err})
		}(addr, idx)
	}
	wg.Wait()
	return
}

// to get domain or ip of url
func urlHandler(urls []string) (validURLs []string) {
	for _, url := range urls {
		s := strings.Replace(url, "https://", "", 1)
		s = strings.Replace(s, "http://", "", 1)
		s = strings.Split(s, "/")[0]
		validURLs = append(validURLs, s)
	}
	return
}

func insertSort(res Results) Results {
	for i := 1; i < len(res); i++ {
		tmp := res[i]
		j := i - 1
		for j >= 0 && tmp.Latency < res[j].Latency {
			res[j+1] = res[j]
			j--
		}
		res[j+1] = tmp
	}

	return res
}