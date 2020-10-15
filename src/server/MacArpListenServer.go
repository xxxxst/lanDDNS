
package server

import (
	"fmt"
	// "bytes"
	"encoding/binary"
	"net"
	// "regexp"
	// "sync"
	"strings"
	// "time"
	
	// "github.com/rjeczalik/notify"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	// . "control"
	. "model"
	// util "util"
)

type MacArpListenServer struct {
	
}

var insMacArpListenServer *MacArpListenServer = nil;

func GetMacArpListenServer() *MacArpListenServer {
	if(insMacArpListenServer == nil) {
		insMacArpListenServer = &MacArpListenServer{};
	}
	return insMacArpListenServer
}

func (c *MacArpListenServer) Run() {
	arr := c.findIfaces();
	if(len(arr) == 0) {
		return;
	}

	// var wg sync.WaitGroup;
	for _,iface := range arr {
		go c.listen(&iface);
		// wg.Add(1);
		// go func(iface net.Interface) {
		// 	defer wg.Done()
		// 	if err := scan(&iface); err != nil {
		// 		log.Printf("interface %v: %v", iface.Name, err)
		// 	}
		// }(iface)
	}

	// go func() {
	// 	wg.Wait();
	// }();
}

func (c *MacArpListenServer) listen(iface *pcap.Interface) {
	fmt.Println(iface.Name);
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Println("aaa", err);
		return;
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		// case <-stop:{
		// 	return;
		// }
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			// if bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
			// 	// This is a packet I sent.
			// 	continue
			// }
			// Note:  we might get some packets here that aren't responses to ones we've sent,
			// if for example someone else sends US an ARP request.  Doesn't much matter, though...
			// all information is good information :)
			fmt.Printf("%v, %v, %v\n", arp.Operation != layers.ARPReply, net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
			continue;
		}
	}
}

func (c *MacArpListenServer) findIfaces() []pcap.Interface {
	md := GetComModel();
	strIp := md.ConfigMd.Server.MacIp;
	arr := strings.Split(strIp, ",");

	// reg := regexp.MustCompile("^([0-9]+\\.[0-9]+\\.[0-9]+).*")

	mapIp := map[uint32] int{};
	for i:=0; i < len(arr); i++ {
		// arr[i] = reg.ReplaceAllString(arr[i], "$1");
		arr[i] = strings.Trim(arr[i], " \t");
		ip := net.ParseIP(arr[i]);
		if(ip == nil) {
			continue;
		}
		ip = ip.To4();
		if(ip == nil) {
			continue;
		}
		numIp := binary.BigEndian.Uint32(ip);
		numIp = numIp & 0xffffff00;
		mapIp[numIp] = 0;
	}
	
	rst := []pcap.Interface{};

	// ifaces, _ := net.Interfaces();
	ifaces, err := pcap.FindAllDevs();
	if(err != nil) {
		return rst;
	}
	// handle err
	for _, i := range ifaces {
		// addrs, _ := i.Addrs()
		addrs := i.Addresses;
		// handle err
		for _, addr := range addrs {
			ip := addr.IP.To4();
			// switch v := addr.(type) {
			// case *net.IPNet:
			// 	// ip = v.IP;
			// 	if !v.IP.IsLoopback() {
			// 		ip4 := v.IP.To4();
			// 		if ip4 != nil {
			// 			//Verify if IP is IPV4
			// 			ip = ip4
			// 		}
			// 	}
			// // case *net.IPAddr:
			// // 	ip = v.IP;
			// }

			if(ip == nil) {
				continue;
			}
			numIp := binary.BigEndian.Uint32(ip);
			numIp = numIp & 0xffffff00;
			_,ok := mapIp[numIp];
			if(ok) {
				// str := ip.String();
				// fmt.Println(str);
				rst = append(rst, i);
				break;
			}

			// process IP address
			// str := ip.String();
			// if str != "<nil>" {
			// 	// for j:= 0; j < len(arr); j++ {
			// 	// 	if(strings.Index(arr[i])) {

			// 	// 	}
			// 	// }
			// 	// arrIp = append(arrIp, str);
			// }
		}
	}

	return rst;
}

func testMacArpListenServer() {
	fmt.Print("");
}