
package server

import (
	"fmt"
	// "bytes"
	"encoding/binary"
	"net"
	"os"
	// "regexp"
	// "sync"
	"strings"
	"time"
	
	// "github.com/rjeczalik/notify"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	. "control"
	. "model"
	// util "util"
)

type IfaceMd struct {
	Iface pcap.Interface;
	NetIp uint32;
	Ip net.IP;
	HardwareAddr net.HardwareAddr;
	Handle *pcap.Handle;
}

type IpData struct {
	Ip uint32;
	Mac string;
}

type MacArpListenServer struct {
	mapAllowIp map[uint32] uint32;
	// mapMacToIp map[string] uint32;
	// mapIpToMac map[uint32] string;
	mapMacIp *OMMap;
	strMacLog string;
	fileMac *os.File;

	chIpData chan IpData;
	mapLicenIfaceName map[string] bool;
}

var insMacArpListenServer *MacArpListenServer = nil;

func GetMacArpListenServer() *MacArpListenServer {
	if(insMacArpListenServer == nil) {
		ins := &MacArpListenServer{};
		ins.mapAllowIp = make(map[uint32] uint32);
		// ins.mapMacToIp = make(map[string] uint32);
		// ins.mapIpToMac = make(map[uint32] string);
		ins.mapMacIp = CreateOMMap();

		ins.chIpData = make(chan IpData, 10);
		ins.mapLicenIfaceName = make(map[string] bool);

		insMacArpListenServer = ins;

	}
	return insMacArpListenServer
}

func (c *MacArpListenServer) Run() {
	c.initConfig();

	if(len(c.mapAllowIp) == 0) {
		return;
	}

	go c.goSetIpMac();

	c.UpdateIface();

	go c.AniUpdateIface();
}

func (c *MacArpListenServer) UpdateIface() {
	arr := c.findIfaces();
	// if(len(arr) == 0) {
	// 	return;
	// }

	for i:=0; i<len(arr); i++ {
		_,ok := c.mapLicenIfaceName[arr[i].Iface.Name];
		if(ok){
			continue;
		}
		c.mapLicenIfaceName[arr[i].Iface.Name] = true;

		md := &arr[i];

		numIp := binary.BigEndian.Uint32(md.Ip);
		c.setIpMac(numIp, md.HardwareAddr.String());

		handle, err := pcap.OpenLive(md.Iface.Name, 65536, true, pcap.BlockForever);
		if err != nil {
			md.Handle = nil;
			continue;
		} else {
			md.Handle = handle;
		}
		go c.listen(*md);
	}

	go c.SendArpData(arr);
}

func (c *MacArpListenServer) AniUpdateIface() {
	ticker := time.NewTimer(time.Second * 120);
	
	for { select {
		case <- ticker.C: {
			c.UpdateIface();
		}
	} }
}

func (c *MacArpListenServer) initConfig() {
	md := GetComModel();
	strIp := md.ConfigMd.Server.MacIp;
	arr := strings.Split(strIp, ",");

	// reg := regexp.MustCompile("^([0-9]+\\.[0-9]+\\.[0-9]+).*")

	mapIp := map[uint32] uint32{};
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

	c.mapAllowIp = mapIp;
}

func (c *MacArpListenServer) getIp(arrIp *[]byte) uint32 {
	if(len(*arrIp) != 4) {
		return 0;
	}
	ip := binary.BigEndian.Uint32(*arrIp);
	return ip;
}

func (c *MacArpListenServer) checkAllowIp(arrIp *[]byte) (bool, uint32) {
	ip := c.getIp(arrIp);
	seg := ip & 0xffffff00;
	_,ok := c.mapAllowIp[seg];
	return ok,ip;
}

// func (c *MacArpListenServer) delIpMac(ip uint32, mac string) {
// 	if val,ok := c.mapMacToIp[mac]; ok {
// 		delete(c.mapIpToMac, val);
// 		delete(c.mapMacToIp, mac);
// 	}
// 	if val, ok := c.mapIpToMac[ip]; ok {
// 		delete(c.mapMacToIp, val);
// 		delete(c.mapIpToMac, ip);
// 	}
// }

func boolCvt(ok bool) string {
	if(ok){
		return "true";
	}else {
		return "false";
	}
}

func (c *MacArpListenServer) goSetIpMac() {
	for {
		select {
		case ipData :=  <- c.chIpData: {
			c.setIpMac(ipData.Ip, ipData.Mac);

			continue;
		}
		}
	}
	// ip := <- chIp;
	// mac := <- chMac;
}

func (c *MacArpListenServer) setIpMac(ip uint32, mac string) {
	a := make(net.IP, 4);
	binary.BigEndian.PutUint32(a, ip);
	strIp := net.IP(a).String();

	// oldIp := c.mapMacIp.GetVal(mac);
	// if(oldIp == nil || oldIp.(uint32) != ip) {
	// 	fmt.Println("save:" + strIp + "," + mac);
	// }
	
	md := GetComModel();

	if(md.ConfigMd.Server.LogMac) {
		oldIp := c.mapMacIp.GetVal(mac);
		if(oldIp != ip) {
			endl := "\r\n";
			str := mac + " " + strIp + endl;
			// c.strMacLog += str;

			if(c.fileMac == nil) {
				c.fileMac, _ = os.OpenFile(md.RootDir + "mac.log", os.O_WRONLY | os.O_APPEND, os.ModePerm);
			}
			if(c.fileMac != nil) {
				c.fileMac.WriteString(str);
			}
		}
	}

	c.mapMacIp.Add(mac, ip);
	
	mainCtl := GetMainCtl();
	mainCtl.DomainMacCtl.SetMacIp(mac, strIp);

	// ipTmp, ok := c.mapMacToIp[mac];
	// if(ok && ipTmp==ip) {
	// 	return;
	// }

	// c.delIpMac(ip, mac);
	// c.mapMacIp.Add(mac, ip);

	// c.mapMacToIp[mac] = ip;
	// c.mapIpToMac[ip] = mac;
}

func (c *MacArpListenServer) checkSaveIp(bIp *[]byte, bMac *[]byte) {
	ok,ip := c.checkAllowIp(bIp);
	if(!ok) {
		return;
	}
	strMac := strings.ToLower(net.HardwareAddr(*bMac).String());
	// c.setIpMac(ip, strMac);

	ipData := IpData{
		Ip: ip,
		Mac: strMac,
	};

	c.chIpData <- ipData;
}

func (c *MacArpListenServer) listen(md IfaceMd) {
	// iface := md.Iface;

	// handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever);
	// if err != nil {
	// 	fmt.Println("aaa", err);
	// 	return;
	// }
	defer md.Handle.Close();

	// fmt.Println(iface.Name, iface.Addresses, iface);

	src := gopacket.NewPacketSource(md.Handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		// case <-stop:{
		// 	return;
		// }
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP);
			if arpLayer == nil {
				continue;
			}
			arp := arpLayer.(*layers.ARP)
			// if bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
			// 	// This is a packet I sent.
			// 	continue
			// }
			// Note:  we might get some packets here that aren't responses to ones we've sent,
			// if for example someone else sends US an ARP request.  Doesn't much matter, though...
			// all information is good information :)
			// ok := false;
			// ip := uint32(0);
			// var arrMac []byte = nil;
			// if(arp.Operation == layers.ARPReply) {
			// 	ok,ip = c.checkAllowIp(&arp.SourceProtAddress);
			// 	arrMac = arp.SourceHwAddress;
			// 	// if(ok) {
			// 	// 	fmt.Printf("-%v, %v, %v\n", arp.Operation != layers.ARPReply, net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress));
			// 	// }
			// } else {
			// 	ok,ip = c.checkAllowIp(&arp.DstProtAddress);
			// 	arrMac = arp.DstHwAddress;
			// 	// if(ok) {
			// 	// 	fmt.Printf(".%v, %v, %v\n", arp.Operation != layers.ARPReply, net.IP(arp.DstProtAddress), net.HardwareAddr(arp.DstHwAddress));
			// 	// }
			// }

			// fmt.Printf("%v, %v, %v\n", arp.Operation == layers.ARPRequest, net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress));
			// fmt.Printf("%v, %v, %v\n", arp.Operation == layers.ARPRequest, net.IP(arp.DstProtAddress), net.HardwareAddr(arp.DstHwAddress));
			// fmt.Printf("----\n");
			
			if(arp.Operation == layers.ARPReply) {
				c.checkSaveIp(&arp.SourceProtAddress, &arp.SourceHwAddress);
				c.checkSaveIp(&arp.DstProtAddress, &arp.DstHwAddress);
			} else if(arp.Operation == layers.ARPRequest) {
				c.checkSaveIp(&arp.SourceProtAddress, &arp.SourceHwAddress);
			}
			// ok,ip := c.checkAllowIp(&arp.SourceProtAddress);
			// if(!ok) {
			// 	continue;
			// }
			// strMac := net.HardwareAddr(arp.SourceHwAddress).String();

			// c.setIpMac(ip, strMac);

			// aaa := make(net.IP, 4);
			// binary.BigEndian.PutUint32(aaa, ip)
			// fmt.Printf("-%v, %v, %v\n", arp.Operation == layers.ARPReply, aaa, strMac);
			continue;
		}
	}
}

func (c *MacArpListenServer) writeARP(md IfaceMd) {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       md.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(md.HardwareAddr),
		SourceProtAddress: []byte(md.Ip),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	numIp := binary.BigEndian.Uint32(md.Ip);

	const maxIpCount int = 255;
	// Send one packet for every address.
	for i:=1; i < maxIpCount; i++ {
		bIp2 := make(net.IP, 4);
		testIp := md.NetIp + uint32(i);
		if(testIp == numIp) {
			continue;
		}
		binary.BigEndian.PutUint32(bIp2, testIp);

		arp.DstProtAddress = bIp2;
		gopacket.SerializeLayers(buf, opts, &eth, &arp)
		if err := md.Handle.WritePacketData(buf.Bytes()); err != nil {
			return;
		}

		time.Sleep(4 * time.Millisecond);
	}
	return;
}

func (c *MacArpListenServer) SendArpData(arr []IfaceMd) {
	time.Sleep(time.Duration(100)*time.Millisecond);
	
	for i:=0; i < len(arr); i++ {
		if(arr[i].Handle == nil) {
			continue;
		}
		c.writeARP(arr[i]);
	}
}

func (c *MacArpListenServer) findIfaces() []IfaceMd {
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
		netIp := numIp & 0xffffff00;
		mapIp[netIp] = 0;
		// fmt.Println("--" + ip.String(), numIp);
	}
	
	rst := []IfaceMd{};

	mapNetIface := make(map[uint32] net.Interface);
	arrNetIfaces, _ := net.Interfaces();
	for i:=0; i < len(arrNetIfaces); i++ {
		netIface := arrNetIfaces[i];

		arr, _ := netIface.Addrs();
		for j:=0; j<len(arr); j++ {
			addr,ok := arr[j].(*net.IPNet);
			if(!ok){
				continue;
			}
			ip := addr.IP.To4();
			if(ip == nil) {
				continue;
			}
			
			numIp := binary.BigEndian.Uint32(ip);
			netIp := numIp & 0xffffff00;
			_,ok = mapIp[netIp];
			// fmt.Println("##" + ip.String(), numIp, arr, netIface);
			if(ok) {
				mapNetIface[numIp] = netIface;
				break;
			}
		}
	}

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
			netIp := numIp & 0xffffff00;
			// fmt.Println("" + ip.String() + "," + i.Name, i.Addresses, numIp, addrs);

			netIface,ok := mapNetIface[numIp];
			if(!ok) {
				continue;
			}

			_,ok = mapIp[netIp];
			if(ok) {
				// str := ip.String();
				// fmt.Println(str);
				md := IfaceMd{
					Iface: i,
					NetIp: netIp,
					HardwareAddr: netIface.HardwareAddr,
					Ip: ip,
				};
				rst = append(rst, md);

				// fmt.Println("" + ip.String() + "," + i.Name, i.Addresses, i);
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