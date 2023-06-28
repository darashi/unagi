package ubloxm8

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// https://qzss.go.jp/en/technical/download/pdf/ps-is-qzss/is-qzss-dcr-010.pdf p.113
var prn2satelliteID = map[byte]byte{
	184: 0x56, // QZS02
	185: 0x57, // QZS04
	189: 0x61, // QZS03
	183: 0x55, // QZS01
	186: 0x58, // QZS1R
}

func ubx2qzqsm(ubx []byte) string {
	if len(ubx) < 14+8*4+2 {
		return ""
	}
	prn := ubx[7] + 182
	satelliteID := prn2satelliteID[prn]

	buf := make([]byte, 32)
	// swap endian
	for i := 0; i < 8; i++ {
		for j := 0; j < 4; j++ {
			buf[i*4+j] = ubx[14+i*4+(3-j)]
		}
	}

	messageType := buf[1] >> 2
	if messageType != 43 && messageType != 44 {
		return ""
	}
	buf[31] &= 0xc0
	buf = buf[:32]
	s := strings.ToUpper(hex.EncodeToString(buf))
	s = fmt.Sprintf("$QZQSM,%02X,%s", satelliteID, s[:len(s)-1])
	s += fmt.Sprintf("*%02X", nmeaChecksum(s))

	return s
}
