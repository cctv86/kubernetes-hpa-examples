package handle

import (
	"log"
	"net"
	"os"

	"github.com/kataras/iris/v12"
)

func GetHostname(ctx iris.Context) {
	host, err := os.Hostname()
	if err != nil {
		log.Printf("%s", err)
	}
	ctx.ContentType("application/json")
	ctx.JSON(host)
}

func GetNameSpace(ctx iris.Context) {
	namespace := getEnv("POD_NAMESPACE", " ")
	ctx.ContentType("application/json")
	ctx.JSON(namespace)
}

func GetIP(ctx iris.Context) {
	ctx.ContentType("application/json")
	ctx.JSON(externalIP())
}

func externalIP() net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip
		}
	}
	return nil
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
