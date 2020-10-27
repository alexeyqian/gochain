package protocol

import (
	"fmt"
)

type IPv4 [4]byte

func NewIPv4(a, b, c, d uint8) IPv4 {
	return IPv4{a, b, c, d}
}

func (ip IPv4) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// four bytes of IPv4 are prepended with 12 other bytes
func (ip IPv4) ToIPv6() []byte {
	return append([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF}, ip[:]...)
}
