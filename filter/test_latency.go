package filter

import (
	"errors"
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

// LatencyTestWithChan test latency between remote and local,
// and store valid remote's index in c
// the faster the less for remote's index
//
// this is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTestWithChan(addrs []string, c chan<- Result, loops int, timeout time.Duration) error {
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
			// By using icmp(ping) to test latency
			p := ping.New(addr)
			p.Timeout = timeout
			p.Count = loops
			err = p.Run()
			statis := p.Statistics()
			if statis.PacketLoss == 100 {
				err = errors.New(addr + ":packet loss reach 100%")
			}
			latency := statis.AvgRtt
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
// for improving efficiency to use LatencyTestWithChan in goroutine instead of this
//
// This is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTest(addrs []string, loops int, timeout time.Duration) (rets []Result, err error) {
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
			p := ping.New(addr)
			p.Timeout = timeout
			p.Count = loops
			err = p.Run()
			statis := p.Statistics()
			if statis.PacketLoss == 100 {
				err = errors.New(addr + ":packet loss reach 100%")
			}
			latency := statis.AvgRtt
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
