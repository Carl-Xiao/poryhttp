package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/lafikl/liblb/p2c"
	"github.com/lafikl/liblb/r2"
)

func TestP2Example(t *testing.T) {
	hosts := []string{"127.0.0.1:8009", "127.0.0.1:8008", "127.0.0.1:8007"}

	// Power of Two choices example
	lb := p2c.New(hosts...)
	for i := 0; i < 10; i++ {
		// uses random power of two choices, because the key length == 0
		host, err := lb.Balance("")
		if err != nil {
			log.Fatal(err)
		}
		// load should be around 33% per host
		fmt.Printf("Send request #%d to host %s\n", i, host)
		// when the work assign to the host is done
		lb.Done(host)
	}
}

//TestNewR2 轮询
func TestNewR2(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	reqPerHost := 100

	lb := r2.New(hosts...)
	loads := map[string]uint64{}

	for i := 0; i < reqPerHost*len(hosts); i++ {
		host, _ := lb.Balance()

		l, _ := loads[host]
		loads[host] = l + 1
	}
	for h, load := range loads {
		if load > uint64(reqPerHost) {
			t.Fatalf("host(%s) got overloaded %d>%d\n", h, load, reqPerHost)
		}
	}
	log.Println(loads)
}

//TestWeightedR2 加权重轮询
func TestWeightedR2(t *testing.T) {
	hosts := []string{"127.0.0.1", "94.0.0.1", "88.0.0.1"}
	reqPerHost := 100

	lb := r2.New()

	// in reverse order just to make sure
	// that insetion order of hosts doesn't affect anything
	for i := len(hosts); i > 0; i-- {
		lb.AddWeight(hosts[i-1], i)
	}

	loads := map[string]uint64{}

	for i := 0; i < reqPerHost*len(hosts); i++ {
		host, _ := lb.Balance()

		l, _ := loads[host]
		loads[host] = l + 1
	}

	for i, host := range hosts {
		fmt.Println(loads[host])

		if loads[host] > uint64(reqPerHost*(i+1)) {
			t.Fatalf("host(%s) got overloaded %d>%d\n", host, loads[host], reqPerHost*i)
		}
	}

}
