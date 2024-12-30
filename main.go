package main

import (
	"C"
	"encoding/binary"
	"time"
	"timemonitor/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)
import (
	"fmt"
	"log"
	"os"
)

const year int64 = 0x83AA7E80

func parseNTPTimeStamp(t layers.NTPTimestamp) time.Time {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t))
	timestamp := int64(binary.BigEndian.Uint32(buf[:4])) - int64(year)
	return time.Unix(timestamp, 0)
}

func main() {
	go utils.Bootstrap()
	s := utils.GetInstance()
	if len(os.Args) != 2 {
		log.Fatal("Wrong param amount\nusage .\\tm [netcard_device_name]\nusage .\\tm -l (list all network card) ")
	}
	var param string = os.Args[1]
	if param == "-l" {
		if d, err := pcap.FindAllDevs(); err != nil {
			log.Fatal(err)
		} else {
			for _, v := range d {
				fmt.Printf("Name:%-55s\tDescription:%-60s\t\n", v.Name, v.Description)
			}
		}
	} else {
		if handle, err := pcap.OpenLive(param, 1600, true, pcap.BlockForever); err != nil {
			log.Fatal(err)
		} else if err := handle.SetBPFFilter("udp and src port 123"); err != nil {
			log.Fatal(err)
		} else {
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for packet := range packetSource.Packets() {
				timestamp := parseNTPTimeStamp(packet.Layer(layers.LayerTypeNTP).(*layers.NTP).TransmitTimestamp)
				ip := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4).DstIP
				s.UpdateLastOKTime(ip, timestamp)
			}
		}
	}
}
