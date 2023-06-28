package ubloxm8

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"

	"go.bug.st/serial"
)

type Receiver struct {
	port serial.Port
}

func NewReceiver(device string, baudrate int) (*Receiver, error) {
	mode := &serial.Mode{BaudRate: baudrate}
	port, err := serial.Open(device, mode)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		port: port,
	}, nil
}

func (r *Receiver) Close() error {
	return r.port.Close()
}

func (r *Receiver) SendUbxCommand(class, id byte, payload []byte) error {
	frame := ubxCommand(class, id, payload)
	_, err := r.port.Write(frame)
	return err
}

func (r *Receiver) EnableQZSSL1S() error {
	return r.SendUbxCommand(0x06, 0x3e, // UBX-CFG-GNSS
		[]byte{
			0x00, // msgVer
			0x20, // numTrkChHw
			0x20, // numTrkChUse
			0x02, // numConfigBlocks = 1
			// config block
			0x00,                   // gnssId = 0 (GPS)
			0x08,                   // resTrkCh
			0x10,                   // maxTrkCh
			0x00,                   // reserved1
			0x01, 0x00, 0x01, 0x01, // flags
			// config block
			0x05,                   // gnssId = 5 (QZSS)
			0x00,                   // resTrkCh
			0x03,                   // maxTrkCh
			0x00,                   // reserved1
			0x01, 0x00, 0x05, 0x05, // flags; enable L1C/A and L1S (in sigCfgMask) and enable=1
		})
}

func (r *Receiver) EnableRXMSFRBXOnURAT1() error {
	return r.SendUbxCommand(0x06, 0x01, // UBX-CFG-MSG; Set message rate(s)
		[]byte{
			0x02, 0x13, // RXM-SFRBX
			0x00, // DDC/I2C
			0x01, // UART 1: set 1
			0x00, // UART 2
			0x00, // USB
			0x00, // SPI
		})
}

const (
	initial = iota
	nmea
	expect_lf
	ubx
)

type MessageHandler func(string)

func (r *Receiver) Receive(onMessage MessageHandler, all bool, verbose bool) error {
	reader := bufio.NewReader(r.port)

	nmeaBuf := ""
	ubxBuf := bytes.Buffer{}
	state := initial
	ubxBytesToRead := 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return err
		}

		switch state {
		case initial:
			bytes, err := reader.Peek(5) // SYNC CHAR 2 (0x62; 1 byte) + CLASS (1 byte) + ID(1 byte) + LENGTH (2 bytes)
			if err != nil {
				return err
			}

			if b == '$' {
				state = nmea
				nmeaBuf = "$"
			} else if b == 0xb5 && bytes[0] == 0x62 {
				state = ubx
				payloadLength := int(binary.LittleEndian.Uint16(bytes[3:5]))
				ubxBytesToRead = 2 + 2 + 2 + payloadLength + 2
				// 1 byte SYNC CHAR 1 (0xB5)
				// 1 byte SYNC CHAR 2 (0x62)
				// 1 byte CLASS
				// 1 byte ID
				// 2 byte LENGTH
				// PAYLOAD
				// 1 byte CK_A
				// 1 byte CK_B
				ubxBuf.Reset()
				ubxBuf.WriteByte(b)
				ubxBytesToRead -= 1 // already read 1 byte (SYNC CHAR 1; 0xB5)
			}
		case nmea:
			if b == '\r' {
				state = expect_lf
			} else {
				nmeaBuf += string(b)
			}
		case expect_lf:
			if b == '\n' {
				state = initial
			} else {
				// unexpected
				log.Println("Warning: LF following CR is missing")
				state = initial
			}
			if all {
				onMessage(nmeaBuf)
			}
		case ubx:
			ubxBuf.WriteByte(b)
			ubxBytesToRead -= 1
			if ubxBytesToRead == 0 {
				state = initial

				buf := ubxBuf.Bytes()
				if ubxBuf.Len() < 8 {
					log.Println("Warning: UBX message too short", buf)
					continue
				}
				// Calculate the UBX checksum from the sequence of CLASS, ID, LENGTH, PAYLOAD
				ck_a, ck_b := ubxChecksum(buf[2 : len(buf)-2])
				if ck_a != buf[len(buf)-2] || ck_b != buf[len(buf)-1] {
					log.Println("Warning: UBX checksum error", buf)
					continue
				}
				if buf[2] == 0x05 && buf[3] == 0x01 {
					if verbose {
						log.Println("UBX-ACK-ACK")
					}
				} else if buf[2] == 0x05 && buf[3] == 0x00 {
					if verbose {
						log.Println("UBX-ACK-NAK")
					}
				} else if buf[2] == 0x02 && buf[3] == 0x13 { // UBX-RXM-SFRBX
					if buf[6] == 0x05 { // QZSS
						qzqsm := ubx2qzqsm(buf)
						if qzqsm != "" {
							onMessage(qzqsm)
						}
					}
				}
			}
		}
	}
}
