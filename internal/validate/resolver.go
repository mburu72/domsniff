package validate

import (
	"context"
	"fmt"
	"net"
	"time"
)

func NewDNSResolver() *net.Resolver {
	servers := []string{"8.8.8.8:53", "1.1.1.1:53"}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			for _, server := range servers {
				d := net.Dialer{Timeout: 5 * time.Second}
				conn, err := d.DialContext(ctx, "udp", server)
				if err == nil {
					return conn, nil
				}
			}
			return nil, fmt.Errorf("all DNS resolvers failed")
		},
	}
}
