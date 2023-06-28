package ubloxm8

import "strings"

func nmeaChecksum(sentence string) byte {
	// NMEA 0183 - Wikipedia
	// https://en.wikipedia.org/wiki/NMEA_0183
	sentence = strings.TrimPrefix(sentence, "$")
	idx := strings.Index(sentence, "*")
	if idx >= 0 {
		sentence = sentence[:idx]
	}

	ck := byte(0)
	for _, b := range []byte(sentence) {
		ck ^= b
	}
	return ck
}
