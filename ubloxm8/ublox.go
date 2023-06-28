package ubloxm8

import (
	"encoding/binary"
)

// https://content.u-blox.com/sites/default/files/products/documents/u-blox8-M8_ReceiverDescrProtSpec_UBX-13003221.pdf

func ubxChecksum(data []byte) (byte, byte) {
	// u-blox 8 / u-blox M8 Receiver description - Manual
	// UBX-13003221 - R28
	// 32.4 UBX Checksum (p. 171)
	ck_a := byte(0)
	ck_b := byte(0)
	for _, b := range data {
		ck_a += b
		ck_b += ck_a
	}

	return ck_a, ck_b
}

func ubxCommand(class byte, id byte, payload []byte) []byte {
	// u-blox 8 / u-blox M8 Receiver description - Manual
	// UBX-13003221 - R28
	// 32.2 UBX Frame Structure
	data := []byte{class, id}
	data = binary.LittleEndian.AppendUint16(data, uint16(len(payload)))
	data = append(data, payload...)
	ck_a, ck_b := ubxChecksum(data)

	frame := []byte{0xb5, 0x62}
	frame = append(frame, data...)
	frame = append(frame, []byte{ck_a, ck_b}...)
	return frame
}
