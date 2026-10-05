//go:build !(darwin || dragonfly || freebsd || netbsd || openbsd)

package listener

import (
	C "github.com/CyberVacation/rostra/constant"
)

func UDPSocketBufferSize() int {
	return C.UDPSocketBufferSize
}
