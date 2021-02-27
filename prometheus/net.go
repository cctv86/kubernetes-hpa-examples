package prometheus

import (
	"github.com/shirou/gopsutil/v3/net"
	"log"
	// "github.com/shirou/gopsutil/mem"  // to use v2
)

func getConnections() float64 {
	v, err := net.Connections("tcp")

	if err != nil {
		log.Println(err)
	}
	i := 0
	for _, value := range v {
		if value.Status == "ESTABLISHED" {
			i++
		}
		continue
	}
	return float64(i)
}
