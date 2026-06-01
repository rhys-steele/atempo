package compose

import (
	"fmt"
	"net"
)

// freePort finds the first available TCP port at or above preferred.
func freePort(preferred int) (int, error) {
	for port := preferred; port < 65535; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		ln.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no free port found starting from %d", preferred)
}
