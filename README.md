# unagi: QZSS DC Report decoder for u-blox M8 receivers


## Overview

This software receives [Satellite Report for Disaster and Crisis Management (DC Report; 災危通報)](https://qzss.go.jp/en/overview/services/sv08_dc-report.html) from [Quasi-Zenith Satellite System (QZSS)](https://qzss.go.jp/en/index.html) with u-blox M8 GPS receivers and converts it to NMEA format.

The code is tested with [Akizuki M-14541 GU-902MGG-USB receiver](https://akizukidenshi.com/catalog/g/gM-14541/), which is equipped with UBX-M8030-KT.

When started, it issues two commands to the receiver. The first is to enable L1S signal reception for QZSS (UBX-CFG-GNSS). The second is to dispatch RXM-SFRBX messages to UART1, which is connected to the PC (UBX-CFG-MSG). After that, it waits for the RXM-SFRBX messages to be received and output them in NMEA format.


## Operation environments

unagi has been tested on the following environments:

* Ubuntu 22.04 LTS (amd64)
* Raspberry Pi OS 64-bit bullseye (Raspberry Pi 4 Model B)

unagi is written in Go and should work on any environment that [Go](https://go.dev/) and [go-serial](https://github.com/bugst/go-serial) work.


## Install

Run the following command:

```
go install github.com/darashi/unagi@latest
```


## Usage

Connect to the receiver to the PC and run the following command:

```
unagi
```

If your receiver is connected to a device other than `/dev/ttyUSB0`, specify the device name with the `-device` option. Details of the options can be found by running `unagi -help`.

You will see the outputs like the following:

```
$QZQSM,56,53ADD371878002B8EA60BA7100000000000000000000000000000012BFA3F94*73
$QZQSM,56,9AADF3710F0002C3E8588ACB118162352C474588FCB1F4165DC00012E112C70*0C
$QZQSM,56,C6ADF3710F0002CC1C5A008B4E2169FB2D400A1F5400000000000013151A7EC*07
$QZQSM,56,53ADD371878002B8EA60BA7100000000000000000000000000000013587C3C8*7E
```

The output is to be sent to the standard output. You can redirect it to a file or pipe it to another program. For example [azarashi](https://github.com/nbtk/azarashi) can be used to obtain human readable information from the output.

If nothing is output, check the environment in which the receiver is located. It needs to be in a location that is open to the sky above. It may be better to try it outdoors. If invoked with the `-all` flag, information other than DC reports will also be output. It may be useful for troubleshooting.


## As a library

You can also use unagi as a library. See `main.go` for the details.


## Other devices

With other devices, even if an M8 series chip is used, it may not work if the internal configuration is different (for example, the code assumes that the PC is connected to UART1 inside the device. If not, then you may be able to make it work with a few changes).


## References

* https://content.u-blox.com/sites/default/files/products/documents/u-blox8-M8_ReceiverDescrProtSpec_UBX-13003221.pdf u-blox 8 / u-blox M8 Receiver description Including protocol specification.
* https://qzss.go.jp/en/technical/download/pdf/ps-is-qzss/is-qzss-dcr-010.pdf Quasi-Zenith Satellite System Interface Specification DC Report Service (IS-QZSS-DCR-010).
* https://eleclog.quitsq.com/2022/12/qzqsm-receiver.html This is an article about creating own receiver using M10, another series of chips from u-blox. The python code was helpful, though the details are different from M8.
* https://twitter.com/Seg_Faul/status/1672963884637093890 My first attempt a few weeks ago didn't go well, but thanks to this tweet it did.
* https://github.com/gpsnmeajp/ub2qzqsm Rust implementation by @Seg_Faul (the author of the above tweet). The code was not publicly available when I was writing unagi, but roughly speaking, unagi is doing almost the same thing.


## License

MIT License
