package filter

import (
	"errors"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/lorenzosaino/go-sysctl"
)

// Result is just result
type Result struct {
	//Name is mirror in eng short
	Name string
	//Address is IP or domain
	Address string
	//Latency in ms
	Latency time.Duration
	//Ping if remote can be ping successfuly
	Ping bool
	//Error for invalid remote host
	Error error
}

// Results for results
type Results []Result

// SortByDuration just sort follow duration
// Using insert sort
func (r Results) SortByDuration() {
	d := sort.Reverse(r)
	sort.Sort(d)
}

// implement sort.Interface interface
func (r Results) Len() int {
	return len(r)
}
func (r Results) Less(i, j int) bool {
	return r[i].Latency > r[j].Latency
}
func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// LatencyTestWithChan test latency between remote and local,
// and store valid remote's index in c
// the faster the less for remote's index
//
// if you do not want a result which sored ,this is better
//
// Args:
// addrs 	 -> []string{"www.bing.com","192.168.12.1"}
// c 	 		 -> store result for goroutine
// loops 	 -> ping times
// gos   	 -> concurrency numbers
// retry   -> wheather try to use net.Dial to test latency after ping failed
// timeout -> timeout for ping and net.Dial
// names 	 -> [optional] url's name exist a map relation
//
// this is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTestWithChan(addrs []string, c chan<- Result, loops int, gos int, retry bool, timeout time.Duration, names ...string) error {
	if len(addrs) < 1 {
		return errors.New("remote host addr must more than one")
	}
	// Check privileged
	if os.Geteuid() != 0 {
		return errors.New("please run this program as root user")
	}
	err := enablePing()
	if err != nil {
		return errors.New("enable ipv4 failed")
	}
	fixedAddrs := urlHandler(addrs)
	var wg sync.WaitGroup
	goroutinePool := make(chan struct{}, gos)
	// To set max concurrency nums
	for i := 0; i < gos; i++ {
		goroutinePool <- struct{}{}
	}
	for idx, addr := range fixedAddrs {
		wg.Add(1)
		<-goroutinePool
		go func(addr string, idx int) {
			defer wg.Done()
			var err error
			var latency time.Duration
			var pingfailed bool
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
				pingfailed = true
				for i := 0; i < loops; i++ {
					_, err = net.DialTimeout("tcp", addr+":http", timeout)
				}
				latency = time.Now().Sub(last) / time.Duration(loops)
			} else {
				latency = statis.AvgRtt
			}
			c <- Result{Name: names[idx], Address: addrs[idx], Latency: latency, Ping: !pingfailed, Error: err}
			goroutinePool <- struct{}{}
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
// gos   -> concurrency numbers
// timeout -> timeout for ping and net.Dial
// names -> [optional] url's name exist a map relation
//
// for improving efficiency to use LatencyTestWithChan in goroutine instead of this
// but this provide a sort method(Results.SortByDuration)
// This is based on icmp protocl, requiring to run as privileged otherwise block
func LatencyTest(addrs []string, loops int, retry bool, gos int, timeout time.Duration, names ...string) (rets Results, err error) {
	if len(addrs) < 1 {
		return nil, errors.New("remote host addr must more than one")
	}
	// Check privileged
	if os.Geteuid() != 0 {
		return nil, errors.New("please run this program as root user")
	}
	err = enablePing()
	if err != nil {
		return nil, errors.New("enable ipv4 failed")
	}
	fixedAddrs := urlHandler(addrs)
	var wg sync.WaitGroup
	goroutinePool := make(chan struct{}, gos)
	// To set max concurrency nums
	for i := 0; i < gos; i++ {
		goroutinePool <- struct{}{}
	}
	for idx, addr := range fixedAddrs {
		wg.Add(1)
		<-goroutinePool
		go func(addr string, idx int) {
			defer wg.Done()
			var err error
			var latency time.Duration
			var pingfailed bool
			p := ping.New(addr)
			p.Timeout = timeout
			p.Count = loops
			err = p.Run()
			statis := p.Statistics()
			// handle remote server banned icmp(ping)
			// switch tp tcp:http
			if statis.PacketLoss == 100 && retry {
				last := time.Now()
				pingfailed = true
				for i := 0; i < loops; i++ {
					_, err = net.DialTimeout("tcp", addr+":http", timeout)
				}
				latency = time.Now().Sub(last) / time.Duration(loops)
			} else {
				latency = statis.AvgRtt
			}
			rets = append(rets, Result{Name: names[idx], Address: addrs[idx], Latency: latency, Ping: !pingfailed, Error: err})
			goroutinePool <- struct{}{}
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

// enable ipv4 ping
// More information @https://pkg.go.dev/github.com/go-ping/ping
func enablePing() error {
	return sysctl.Set("net.ipv4.ping_group_range", "0 2147483647")
}
