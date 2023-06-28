package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/darashi/unagi/ubloxm8"
)

var (
	flagDevice   = flag.String("device", "/dev/ttyUSB0", "serial port")
	flagBaudRate = flag.Int("baud-rate", 9600, "baud rate")
	flagAll      = flag.Bool("all", false, "output all messages; not only $QZQSM")
	flagVerbose  = flag.Bool("verbose", false, "verbose output")
)

func main() {
	flag.Parse()

	receiver, err := ubloxm8.NewReceiver(*flagDevice, *flagBaudRate)
	if err != nil {
		log.Fatal(err)
	}
	defer receiver.Close()

	if err := receiver.EnableQZSSL1S(); err != nil {
		log.Fatal(err)
	}

	if err := receiver.EnableRXMSFRBXOnURAT1(); err != nil {
		log.Fatal(err)
	}

	handler := func(msg string) {
		fmt.Println(msg)
	}

	if err := receiver.Receive(handler, *flagAll, *flagVerbose); err != nil {
		log.Fatal(err)
	}
}
